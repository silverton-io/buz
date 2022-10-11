// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputcloudevents

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/protocol"
	"github.com/tidwall/gjson"
)

func BuildEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []envelope.Envelope {
	var envelopes []envelope.Envelope
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not read request body")
		return envelopes
	}
	for _, ce := range gjson.ParseBytes(reqBody).Array() {
		n := envelope.BuildCommonEnvelope(c, conf.Middleware, m)
		cEvent, err := buildEvent(ce)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not build Cloudevent")
		}
		// Event Meta
		n.EventMeta.Protocol = protocol.CLOUDEVENTS
		n.EventMeta.Schema = cEvent.DataSchema
		// Source
		if cEvent.Time != nil {
			n.Pipeline.Source.GeneratedTstamp = cEvent.Time
			n.Pipeline.Source.SentTstamp = cEvent.Time
		}
		// Payload
		n.Payload = cEvent.Data
		envelopes = append(envelopes, n)
	}
	return envelopes
}
