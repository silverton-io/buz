package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func CloudeventsHandler(h EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.ContentType() == "application/cloudevents+json" || c.ContentType() == "application/cloudevents-batch+json" {
			envelopes := envelope.BuildCloudeventEnvelopesFromRequest(c, *h.Config)
			annotatedEnvelopes := validator.Annotate(envelopes, h.Cache)
			h.Manifold.Enqueue(annotatedEnvelopes)
			c.JSON(http.StatusOK, response.Ok)
		} else {
			c.JSON(http.StatusBadRequest, response.InvalidContentType)
		}
	}
	return gin.HandlerFunc(fn)
}
