// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputsnowplow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/params"
	"github.com/silverton-io/buz/pkg/response"
)

func Handler(h params.Handler) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		err := h.Manifold.Distribute(envelopes)
		if err != nil {
			c.Header("Retry-After", response.RETRY_AFTER_60)
			c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
		} else {
			c.JSON(http.StatusOK, response.Ok)
		}
		if c.Request.Method == http.MethodGet {
			redirectUrl, _ := c.GetQuery("u")
			if redirectUrl != "" && h.Config.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("🟢 redirecting to " + redirectUrl)
				c.Redirect(http.StatusFound, redirectUrl)
			}
		}
	}
	return gin.HandlerFunc(fn)
}
