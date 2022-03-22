package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/response"
)

func WebhookHandler(handlerParams EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.ContentType() == "application/json" {
			envelopes := envelope.BuildWebhookEnvelopesFromRequest(c)
			handlerParams.Manifold.Enqueue(envelopes)
			c.JSON(http.StatusOK, response.Ok)
		} else {
			c.JSON(http.StatusBadRequest, response.InvalidContentType)
		}
	}
	return gin.HandlerFunc(fn)
}
