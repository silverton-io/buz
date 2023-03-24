// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package snowplow

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/protocol"
	"github.com/silverton-io/buz/pkg/util"
	"github.com/tidwall/gjson"
)

func buildSnowplowEnvelope(e SnowplowEvent) envelope.Envelope {
	n := envelope.NewEnvelope()
	n.Timestamp = *e.DvceCreatedTstamp
	n.Protocol = protocol.SNOWPLOW
	n.Schema = *e.SelfDescribingEvent.SchemaName()
	n.Contexts = e.Contexts
	n.Payload = e.Map()
	return n
}

func buildEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []envelope.Envelope {
	var envelopes []envelope.Envelope
	if c.Request.Method == "POST" {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not read request body")
			return envelopes
		}
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			spEvent := buildEventFromMappedParams(c, event.Value().(map[string]interface{}), *conf)
			e := buildSnowplowEnvelope(spEvent)
			envelopes = append(envelopes, e)
		}
	} else {
		params := util.MapUrlParams(c)
		spEvent := buildEventFromMappedParams(c, params, *conf)
		e := buildSnowplowEnvelope(spEvent)
		envelopes = append(envelopes, e)
	}
	return envelopes
}
