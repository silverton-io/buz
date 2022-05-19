package envelope

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/generic"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/pixel"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/webhook"
	"github.com/tidwall/gjson"
)

func buildRelayEventFromPayload(proto string, payload []byte) event.Event {
	switch proto {
	case protocol.SNOWPLOW:
		e := event.SelfDescribingPayload{}
		json.Unmarshal(payload, &e)
		return e
	case protocol.CLOUDEVENTS:
		e := cloudevents.CloudEvent{}
		json.Unmarshal(payload, &e)
		return e
	case protocol.GENERIC:
		e := generic.GenericEvent{}
		json.Unmarshal(payload, &e)
		return e
	case protocol.WEBHOOK:
		e := webhook.WebhookEvent{}
		json.Unmarshal(payload, &e)
		return e
	case protocol.PIXEL:
		e := pixel.PixelEvent{}
		json.Unmarshal(payload, &e)
		return e
	}
	return nil
}

func BuildRelayEnvelopesFromRequest(c *gin.Context, m *meta.CollectorMeta) []Envelope {
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
		protocol := relayedEvent.Get("eventMeta.protocol").String()
		payload := relayedEvent.Get("payload").Raw
		e := buildRelayEventFromPayload(protocol, []byte(payload))
		var n Envelope
		json.Unmarshal([]byte(relayedEvent.Raw), &n)
		t := time.Now().UTC()
		id := uuid.New()
		// Relay Meta
		n.Pipeline.Relay.Relayed = true
		n.Pipeline.Relay.Id = &id
		n.Pipeline.Relay.Tstamp = &t
		// Payload
		n.Payload = e
		envelopes = append(envelopes, n)
	}
	return envelopes
}
