package envelope

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/generic"
	"github.com/silverton-io/honeypot/pkg/pixel"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/silverton-io/honeypot/pkg/webhook"
	"github.com/tidwall/gjson"
)

func buildSnowplowEnvelope(c *gin.Context, e snowplow.SnowplowEvent) Envelope {
	envelope := Envelope{
		SourceMetadata: SourceMetadata{
			Ip: c.ClientIP(),
		},
		EventMetadata: EventMetadata{
			Uuid:     uuid.New(),
			Protocol: protocol.SNOWPLOW,
		},
		CollectorMetadata: CollectorMetadata{
			Tstamp: time.Now().UTC(),
		},
		RelayMetadata: RelayMetadata{
			IsRelayed: false,
		},
		ValidationMetadata: ValidationMetadata{},
		Payload:            e,
	}
	return envelope
}

func BuildSnowplowEnvelopesFromRequest(c *gin.Context, conf config.Config) []Envelope {
	var envelopes []Envelope
	if c.Request.Method == "POST" {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not read request body")
			return envelopes
		}
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			spEvent := snowplow.BuildEventFromMappedParams(c, event.Value().(map[string]interface{}), conf)
			e := buildSnowplowEnvelope(c, spEvent)
			envelopes = append(envelopes, e)
		}
	} else {
		params := util.MapUrlParams(c)
		spEvent := snowplow.BuildEventFromMappedParams(c, params, conf)
		e := buildSnowplowEnvelope(c, spEvent)
		envelopes = append(envelopes, e)
	}
	return envelopes
}

func BuildGenericEnvelopesFromRequest(c *gin.Context, conf config.Config) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		genEvent := generic.BuildEvent(e, conf.Generic)
		envelope := Envelope{
			SourceMetadata: SourceMetadata{Ip: c.ClientIP()},
			EventMetadata: EventMetadata{
				Uuid:     uuid.New(),
				Protocol: protocol.GENERIC,
			},
			CollectorMetadata: CollectorMetadata{
				Tstamp: time.Now().UTC(),
			},
			RelayMetadata: RelayMetadata{
				IsRelayed: false,
			},
			ValidationMetadata: ValidationMetadata{},
			Payload:            genEvent,
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
		return envelopes
	}
	for _, ce := range gjson.ParseBytes(reqBody).Array() {
		cEvent, _ := cloudevents.BuildEvent(ce)
		envelope := Envelope{
			SourceMetadata: SourceMetadata{
				Ip: c.ClientIP(),
			},
			EventMetadata: EventMetadata{
				Uuid:     uuid.New(),
				Protocol: protocol.CLOUDEVENTS,
			},
			CollectorMetadata: CollectorMetadata{
				Tstamp: time.Now().UTC(),
			},
			RelayMetadata: RelayMetadata{
				IsRelayed: false,
			},
			ValidationMetadata: ValidationMetadata{},
			Payload:            cEvent,
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func BuildWebhookEnvelopesFromRequest(c *gin.Context, conf config.Config) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		whEvent, err := webhook.BuildEvent(c, e)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not build WebhookEvent")
		}
		envelope := Envelope{
			SourceMetadata: SourceMetadata{
				Ip: c.ClientIP(),
			},
			EventMetadata: EventMetadata{
				Uuid:     uuid.New(),
				Protocol: protocol.WEBHOOK,
			},
			CollectorMetadata: CollectorMetadata{
				Tstamp: time.Now().UTC(),
			},
			RelayMetadata: RelayMetadata{
				IsRelayed: false,
			},
			ValidationMetadata: ValidationMetadata{
				IsValid: true,
			},
			Payload: whEvent,
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func BuildPixelEnvelopesFromRequest(c *gin.Context, conf config.Config) []Envelope {
	var envelopes []Envelope
	params := util.MapUrlParams(c)
	var urlNames = make(map[string]string)
	for _, i := range conf.Pixel.Paths {
		urlNames[i.Path] = i.Name
	}
	name := urlNames[c.Request.URL.Path]
	pEvent, err := pixel.BuildEvent(c, name, params)
	if err != nil {
		log.Error().Err(err).Msg("could not build PixelEvent")
	}
	envelope := Envelope{
		SourceMetadata: SourceMetadata{
			Ip: c.ClientIP(),
		},
		EventMetadata: EventMetadata{
			Uuid:     uuid.New(),
			Protocol: protocol.PIXEL,
		},
		CollectorMetadata: CollectorMetadata{
			Tstamp: time.Now().UTC(),
		},
		RelayMetadata: RelayMetadata{
			IsRelayed: false,
		},
		ValidationMetadata: ValidationMetadata{
			IsValid: true,
		},
		Payload: pEvent,
	}
	envelopes = append(envelopes, envelope)
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
		case protocol.WEBHOOK:
			payload := webhook.WebhookEvent{}
			json.Unmarshal([]byte(eventPayload), &payload)
			envelope.Payload = payload
		default:
			payload := snowplow.SnowplowEvent{}
			json.Unmarshal([]byte(eventPayload), &payload)
			envelope.Payload = payload
		}
		relayMeta := RelayMetadata{
			IsRelayed: true,
			RelayId:   uuid.New(),
		}
		envelope.RelayMetadata = relayMeta
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}
