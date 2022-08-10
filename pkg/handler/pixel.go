// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/annotator"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/params"
	"github.com/silverton-io/honeypot/pkg/privacy"
	"github.com/silverton-io/honeypot/pkg/response"
)

const PX string = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8Xw8AAoMBgDTD2qgAAAAASUVORK5CYII="

func PixelHandler(h params.Handler) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := envelope.BuildPixelEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		annotatedEnvelopes := annotator.Annotate(envelopes, h.Cache)
		anonymizedEnvelopes := privacy.AnonymizeEnvelopes(annotatedEnvelopes, h.Config.Privacy)
		err := h.Manifold.Distribute(anonymizedEnvelopes, h.ProtocolStats)
		if err != nil {
			c.Header("Retry-After", response.RETRY_AFTER_60)
			c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
		} else {
			b, _ := base64.StdEncoding.DecodeString(PX)
			c.Data(http.StatusOK, "image/png", b)
		}
	}
	return gin.HandlerFunc(fn)
}
