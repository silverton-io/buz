// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputSelfDescribing

import (
	"bytes"
	"compress/gzip"
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
	// If the request body is gzipped, decompress it
	if c.GetHeader("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(bytes.NewReader(reqBody))
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not decompress gzipped request body")
			return envelopes
		}
		defer reader.Close()
		reqBody, err = io.ReadAll(reader)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not read decompressed gzipped request body")
			return envelopes
		}
	}

	for _, e := range gjson.ParseBytes(reqBody).Array() {
		n := envelope.BuildCommonEnvelope(c, conf.Middleware, m)
		genEvent, err := buildEvent(e, conf.SelfDescribing)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not build generic event")
		}
		// Event meta
		n.EventMeta.Protocol = protocol.SELF_DESCRIBING
		n.EventMeta.Schema = genEvent.Payload.Schema
		// Context
		n.Contexts = &genEvent.Contexts
		// Payload
		n.Payload = genEvent.Payload.Data
		envelopes = append(envelopes, n)
	}
	return envelopes
}
