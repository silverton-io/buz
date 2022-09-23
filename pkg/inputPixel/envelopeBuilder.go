// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputpixel

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/protocol"
)

// NOTE - one envelope per request
func BuildEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []envelope.Envelope {
	var envelopes []envelope.Envelope
	sde, err := buildEvent(c)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not build pixel event")
	}
	contexts := envelope.BuildContextsFromRequest(c)
	n := envelope.BuildCommonEnvelope(c, conf.Middleware, m)
	// Event Meta
	n.EventMeta.Protocol = protocol.PIXEL
	n.EventMeta.Schema = sde.Schema
	// Contexts
	n.Contexts = &contexts
	// Payload
	n.Payload = sde.Data
	envelopes = append(envelopes, n)
	return envelopes
}
