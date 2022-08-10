// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

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
	sde, err := pixel.BuildEvent(c)
	if err != nil {
		log.Error().Err(err).Msg("could not build pixel event")
	}
	contexts := buildContextsFromRequest(c)
	n := buildCommonEnvelope(c, conf.Middleware, m)
	// Event Meta
	n.EventMeta.Protocol = protocol.PIXEL
	n.EventMeta.Schema = sde.Schema
	// Contexts
	n.Contexts = contexts
	// Payload
	n.Payload = sde.Data
	envelopes = append(envelopes, n)
	return envelopes
}
