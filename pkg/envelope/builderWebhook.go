package envelope

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/webhook"
	"github.com/tidwall/gjson"
)

func BuildWebhookEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		n := buildCommonEnvelope(c, m)
		contexts := buildContextFromRequest(c)
		whEvent, err := webhook.BuildEvent(c, e)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not build webhook event")
		}
		// Event Meta
		n.EventMeta.Protocol = protocol.WEBHOOK
		// Contexts
		n.Contexts = contexts
		// Payload
		n.Payload = whEvent
		envelopes = append(envelopes, n)
	}
	return envelopes
}
