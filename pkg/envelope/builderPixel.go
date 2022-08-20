// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/pixel"
	"github.com/silverton-io/buz/pkg/protocol"
)

// NOTE - one envelope per request
func BuildPixelEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
	var envelopes []Envelope
	sde, err := pixel.BuildEvent(c)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not build pixel event")
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
