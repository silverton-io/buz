package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func GenericHandler(h EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := envelope.BuildGenericEnvelopesFromRequest(c, *h.Config)
		annotatedEnvelopes := validator.Annotate(envelopes, h.Cache)
		h.Manifold.Enqueue(annotatedEnvelopes)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
