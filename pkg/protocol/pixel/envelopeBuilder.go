// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package pixel

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/protocol"
)

// NOTE - one envelope per request
func buildEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []envelope.Envelope {
	var envelopes []envelope.Envelope
	contexts := envelope.BuildContextsFromRequest(c)
	n := envelope.NewEnvelope(conf.App)
	evnt, err := buildEvent(c)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not build pixel event")
	}
	n.Protocol = protocol.PIXEL
	if evnt.Schema != "" {
		n.Schema = evnt.Schema
	}
	n.Contexts = &contexts
	n.Payload = evnt.Data
	envelopes = append(envelopes, n)
	return envelopes
}
