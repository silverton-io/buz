package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func SnowplowHandler(h EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := envelope.BuildSnowplowEnvelopesFromRequest(c, *h.Config)
		validEnvelopes, invalidEnvelopes, stats := validator.BifurcateAndAnnotate(envelopes, h.Cache)
		h.Sink.BatchPublishValidAndInvalid(ctx, validEnvelopes, invalidEnvelopes)
		h.Meta.ProtocolStats.Merge(&stats)
		if c.Request.Method == http.MethodGet {
			redirectUrl, _ := c.GetQuery("u")
			if redirectUrl != "" && h.Config.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("redirecting to " + redirectUrl)
				c.Redirect(http.StatusFound, redirectUrl)
			}
		}
		c.JSON(http.StatusOK, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
