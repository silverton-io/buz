package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/annotator"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/response"
)

func GenericHandler(h EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.ContentType() == "application/json" {
			envelopes := envelope.BuildGenericEnvelopesFromRequest(c, *h.Config)
			annotatedEnvelopes := annotator.Annotate(envelopes, h.Cache)
			err := h.Manifold.Distribute(annotatedEnvelopes)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, response.DistributionError)
			} else {
				c.JSON(200, response.Ok)
			}

		} else {
			c.JSON(http.StatusBadRequest, response.InvalidContentType)
		}
	}
	return gin.HandlerFunc(fn)
}
