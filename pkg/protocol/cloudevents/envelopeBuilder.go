// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package cloudevents

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

func newCloudeventsEnvelope(event CloudEvent) envelope.Envelope {
	n := envelope.NewEnvelope()
	// Event Meta
	n.Protocol = protocol.CLOUDEVENTS
	n.Schema = event.DataSchema
	// Set timestamp according to upstream source
	if event.Time != nil {
		n.Timestamp = *event.Time
	}
	// Payload
	n.Payload = event.Data
	return n
}

func buildEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []envelope.Envelope {
	var envelopes []envelope.Envelope
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not read request body")
		return envelopes
	}
	for _, ce := range gjson.ParseBytes(reqBody).Array() {
		cEvent, err := buildEvent(ce)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not build Cloudevent")
		}
		n := newCloudeventsEnvelope(cEvent)
		envelopes = append(envelopes, n)
	}
	return envelopes
}
