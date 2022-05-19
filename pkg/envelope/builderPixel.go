package envelope

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/pixel"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/util"
)

func BuildPixelEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
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
	n := buildCommonEnvelope(c, m)
	// Event Meta
	n.EventMeta.Protocol = protocol.PIXEL
	// Payload
	n.Payload = pEvent
	envelopes = append(envelopes, n)
	return envelopes
}
