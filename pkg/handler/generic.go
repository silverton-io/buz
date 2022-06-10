package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/annotator"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/params"
	"github.com/silverton-io/honeypot/pkg/privacy"
	"github.com/silverton-io/honeypot/pkg/response"
)

func GenericHandler(h params.Handler) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.ContentType() == "application/json" {
			envelopes := envelope.BuildGenericEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
			annotatedEnvelopes := annotator.Annotate(envelopes, h.Cache)
			anonymizedEnvelopes := privacy.AnonymizeEnvelopes(annotatedEnvelopes, h.Config.Privacy)
			err := h.Manifold.Distribute(anonymizedEnvelopes, h.ProtocolStats)
			if err != nil {
				c.Header("Retry-After", response.RETRY_AFTER_60)
				c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
			} else {
				c.JSON(200, response.Ok)
			}

		} else {
			c.JSON(http.StatusBadRequest, response.InvalidContentType)
		}
	}
	return gin.HandlerFunc(fn)
}
