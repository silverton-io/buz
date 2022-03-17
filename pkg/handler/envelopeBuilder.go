package handler

import (
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	e "github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/generic"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/request"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/tidwall/gjson"
)

func buildSnowplowEnvelopesFromRequest(c *gin.Context, conf config.Config) []e.Envelope {
	var envelopes []e.Envelope
	isRelayed := false
	if c.Request.Method == "POST" {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not read request body")
		}
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			spEvent := snowplow.BuildEventFromMappedParams(c, event.Value().(map[string]interface{}), conf)
			schema := spEvent.Schema()
			envelope := e.Envelope{
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
		envelope := e.Envelope{
			EventProtocol: protocol.SNOWPLOW,
			EventSchema:   schema,
			Tstamp:        time.Now(),
			Ip:            *spEvent.User_ipaddress,
			Payload:       spEvent,
			IsRelayed:     &isRelayed,
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func buildGenericEnvelopesFromRequest(c *gin.Context, conf config.Config) []e.Envelope {
	var envelopes []e.Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
	}
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		genEvent := generic.BuildEvent(e, conf.Generic)
		isRelayed := false
		envelope := event.Envelope{
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

func buildCloudeventEnvelopesFromRequest(c *gin.Context, conf config.Config) []e.Envelope {
	var envelopes []e.Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
	}
	for _, ce := range gjson.ParseBytes(reqBody).Array() {
		cEvent, _ := cloudevents.BuildEvent(ce)
		isRelayed := false
		envelope := e.Envelope{
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
