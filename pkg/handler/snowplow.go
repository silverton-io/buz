package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/annotator"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/params"
	"github.com/silverton-io/honeypot/pkg/privacy"
	"github.com/silverton-io/honeypot/pkg/response"
)

func SnowplowHandler(h params.Handler) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := envelope.BuildSnowplowEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		annotatedEnvelopes := annotator.Annotate(envelopes, h.Cache)
		anonymizedEnvelopes := privacy.AnonymizeEnvelopes(annotatedEnvelopes, h.Config.Privacy)
		err := h.Manifold.Distribute(anonymizedEnvelopes, h.ProtocolStats)
		if err != nil {
			c.Header("Retry-After", response.RETRY_AFTER_60)
			c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
		} else {
			c.JSON(http.StatusOK, response.Ok)
		}
		if c.Request.Method == http.MethodGet {
			redirectUrl, _ := c.GetQuery("u")
			if redirectUrl != "" && h.Config.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("redirecting to " + redirectUrl)
				c.Redirect(http.StatusFound, redirectUrl)
			}
		}
	}
	return gin.HandlerFunc(fn)
}
