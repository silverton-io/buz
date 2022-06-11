package envelope

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/pixel"
	"github.com/silverton-io/honeypot/pkg/protocol"
)

// NOTE - one envelope per request
func BuildPixelEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
	var envelopes []Envelope
	pEvent, err := pixel.BuildEvent(c)
	if err != nil {
		log.Error().Err(err).Msg("could not build pixel event")
	}
	contexts := buildContextsFromRequest(c)
	n := buildCommonEnvelope(c, m)
	// Event Meta
	n.EventMeta.Protocol = protocol.PIXEL
	// Contexts
	n.Contexts = contexts
	// Payload
	n.Payload = pEvent
	envelopes = append(envelopes, n)
	return envelopes
}
