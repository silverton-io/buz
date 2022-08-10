// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/annotator"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/params"
	"github.com/silverton-io/honeypot/pkg/protocol"
)

func SquawkboxHandler(h params.Handler, eventProtocol string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var envelopes []envelope.Envelope
		switch eventProtocol {
		case protocol.SNOWPLOW:
			envelopes = envelope.BuildSnowplowEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		case protocol.CLOUDEVENTS:
			envelopes = envelope.BuildCloudeventEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		case protocol.GENERIC:
			envelopes = envelope.BuildGenericEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		}
		annotatedEnvelopes := annotator.Annotate(envelopes, h.Cache)
		c.JSON(http.StatusOK, annotatedEnvelopes)
	}
	return gin.HandlerFunc(fn)
}
