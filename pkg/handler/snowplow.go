package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/annotator"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/response"
)

func SnowplowHandler(h EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := envelope.BuildSnowplowEnvelopesFromRequest(c, *h.Config)
		annotatedEnvelopes := annotator.Annotate(envelopes, h.Cache)
		err := h.Manifold.Distribute(annotatedEnvelopes)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, response.DistributionError)
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
