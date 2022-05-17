package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/annotator"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/params"
	"github.com/silverton-io/honeypot/pkg/response"
)

func CloudeventsHandler(h params.Handler) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.ContentType() == "application/cloudevents+json" || c.ContentType() == "application/cloudevents-batch+json" {
			envelopes := envelope.BuildCloudeventEnvelopesFromRequest(c, *h.Config)
			annotatedEnvelopes := annotator.Annotate(envelopes, h.Cache)
			err := h.Manifold.Distribute(annotatedEnvelopes, h.Meta)
			if err != nil {
				c.Header("Retry-After", response.RETRY_AFTER_60)
				c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
			} else {
				c.JSON(http.StatusOK, response.Ok)
			}
		} else {
			c.JSON(http.StatusBadRequest, response.InvalidContentType)
		}
	}
	return gin.HandlerFunc(fn)
}
