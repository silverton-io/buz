// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package webhook

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

func buildEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []envelope.Envelope {
	var envelopes []envelope.Envelope
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not read request body")
		return envelopes
	}
	// If the request body is gzipped, decompress it
    if c.GetHeader("Content-Encoding") == "gzip" {
        log.Debug().Msg("🟢 Request body is gzip encoded - attempting to decompress.")
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
		n := envelope.NewEnvelope(conf.App)
		contexts := envelope.BuildContextsFromRequest(c)
		evnt, err := buildEvent(c, e)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not build webhook event")
		}
		n.Protocol = protocol.WEBHOOK
		if evnt.Schema != "" {
			n.Schema = evnt.Schema
		}
		n.Contexts = &contexts
		n.Payload = evnt.Data
		envelopes = append(envelopes, n)
	}
	return envelopes
}
