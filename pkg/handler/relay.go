package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/generic"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func buildRelayEnvelopesFromRequest(c *gin.Context) []event.Envelope {
	var envelopes []event.Envelope
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	relayedEvents := gjson.ParseBytes(reqBody)
	for _, relayedEvent := range relayedEvents.Array() {
		eventProtocol := relayedEvent.Get("eventProtocol").String()
		eventPayload := relayedEvent.Get("payload").Raw
		var envelope event.Envelope
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

// RelayHandler processes incoming envelopes, splits them in half,
// and sends them to the configured sink. It relies on upstream validation.
func RelayHandler(p EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := buildRelayEnvelopesFromRequest(c)
		validEnvelopes, invalidEnvelopes := validator.Bifurcate(envelopes)
		p.Sink.BatchPublishValidAndInvalid(ctx, protocol.RELAY, validEnvelopes, invalidEnvelopes, p.Meta)
	}
	return gin.HandlerFunc(fn)
}
