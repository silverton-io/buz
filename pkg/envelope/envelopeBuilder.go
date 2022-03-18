package envelope

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/generic"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/request"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/silverton-io/honeypot/pkg/webhook"
	"github.com/tidwall/gjson"
)

func BuildSnowplowEnvelopesFromRequest(c *gin.Context, conf config.Config) []Envelope {
	var envelopes []Envelope
	isRelayed := false
	if c.Request.Method == "POST" {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not read request body")
			envelope := Envelope{}
			envelopes = append(envelopes, envelope)
			return envelopes
		}
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			spEvent := snowplow.BuildEventFromMappedParams(c, event.Value().(map[string]interface{}), conf)
			schema := spEvent.Schema()
			envelope := Envelope{
				EventProtocol: protocol.SNOWPLOW,
				EventSchema:   schema,
				Tstamp:        time.Now(),
				Ip:            *spEvent.User_ipaddress,
				Payload:       spEvent,
				IsRelayed:     &isRelayed,
			}
			envelopes = append(envelopes, envelope)
		}
	} else {
		params := request.MapParams(c)
		spEvent := snowplow.BuildEventFromMappedParams(c, params, conf)
		schema := spEvent.Schema()
		e := Envelope{
			EventProtocol: protocol.SNOWPLOW,
			EventSchema:   schema,
			Tstamp:        time.Now(),
			Ip:            *spEvent.User_ipaddress,
			Payload:       spEvent,
			IsRelayed:     &isRelayed,
		}
		envelopes = append(envelopes, e)
	}
	return envelopes
}

func BuildGenericEnvelopesFromRequest(c *gin.Context, conf config.Config) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		envelope := Envelope{}
		envelopes = append(envelopes, envelope)
		return envelopes
	}
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		genEvent := generic.BuildEvent(e, conf.Generic)
		isRelayed := false
		envelope := Envelope{
			EventProtocol: protocol.GENERIC,
			EventSchema:   &genEvent.Payload.Schema,
			Tstamp:        time.Now(),
			Ip:            c.ClientIP(),
			Payload:       genEvent,
			IsRelayed:     &isRelayed,
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func BuildCloudeventEnvelopesFromRequest(c *gin.Context, conf config.Config) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		envelope := Envelope{}
		envelopes = append(envelopes, envelope)
		return envelopes
	}
	for _, ce := range gjson.ParseBytes(reqBody).Array() {
		cEvent, _ := cloudevents.BuildEvent(ce)
		isRelayed := false
		envelope := Envelope{
			EventProtocol: protocol.CLOUDEVENTS,
			EventSchema:   &cEvent.DataSchema,
			Tstamp:        time.Now(),
			Source:        cEvent.Source,
			Ip:            c.ClientIP(),
			Payload:       cEvent,
			IsRelayed:     &isRelayed,
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func BuildRelayEnvelopesFromRequest(c *gin.Context) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		envelope := Envelope{}
		envelopes = append(envelopes, envelope)
		return envelopes
	}
	relayedEvents := gjson.ParseBytes(reqBody)
	for _, relayedEvent := range relayedEvents.Array() {
		eventProtocol := relayedEvent.Get("eventProtocol").String()
		eventPayload := relayedEvent.Get("payload").Raw
		var envelope Envelope
		json.Unmarshal([]byte(relayedEvent.Raw), &envelope)
		switch eventProtocol {
		case protocol.SNOWPLOW:
			payload := snowplow.SnowplowEvent{}
			json.Unmarshal([]byte(eventPayload), &payload)
			envelope.Payload = payload
		case protocol.CLOUDEVENTS:
			payload := cloudevents.CloudEvent{}
			json.Unmarshal([]byte(eventPayload), &payload)
			envelope.Payload = payload
		case protocol.GENERIC:
			payload := generic.GenericEvent{}
			json.Unmarshal([]byte(eventPayload), &payload)
			envelope.Payload = payload
		default:
			payload := snowplow.SnowplowEvent{}
			json.Unmarshal([]byte(eventPayload), &payload)
			envelope.Payload = payload
		}
		isRelayed := true
		envelope.IsRelayed = &isRelayed
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func BuildWebhookEnvelopesFromRequest(c *gin.Context) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		envelope := Envelope{}
		envelopes = append(envelopes, envelope)
		return envelopes
	}
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		whEvent, err := webhook.BuildEvent(e)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not build WebhookEvent")
		}
		isRelayed := false
		envelope := Envelope{
			EventProtocol: protocol.WEBHOOK,
			EventSchema:   whEvent.Schema(),
			Tstamp:        time.Now(),
			Ip:            c.ClientIP(),
			Payload:       whEvent,
			IsRelayed:     &isRelayed,
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}
