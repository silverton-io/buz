package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/generic"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func RelayHandler(p EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var envelopes []event.Envelope
		ctx := context.Background()
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		relayedEvents := gjson.ParseBytes(reqBody)
		for _, relayedEvent := range relayedEvents.Array() {
			eventProtocol := relayedEvent.Get("eventProtocol").String()
			eventPayload := relayedEvent.Get("payload").Raw
			var envelope event.Envelope
			envelopeErr := json.Unmarshal([]byte(relayedEvent.Raw), &envelope)
			if envelopeErr != nil {
				log.Error().Stack().Err(envelopeErr).Msg("could not unmarshal relayed event to envelope")
			}
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
			util.Pprint(envelope)
		}
		validEnvelopes, invalidEnvelopes := validator.Bifurcate(envelopes)
		p.Sink.BatchPublishValidAndInvalid(ctx, protocol.RELAY, validEnvelopes, invalidEnvelopes, p.Meta)
	}
	return gin.HandlerFunc(fn)
}
