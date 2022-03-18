package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func CloudeventsHandler(handlerParams EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.ContentType() == "application/cloudevents+json" || c.ContentType() == "application/cloudevents-batch+json" {
			ctx := context.Background()
			envelopes := envelope.BuildCloudeventEnvelopesFromRequest(c, *handlerParams.Config)
			validEvents, invalidEvents := validator.BifurcateAndAnnotate(envelopes, handlerParams.Cache)
			handlerParams.Sink.BatchPublishValidAndInvalid(ctx, protocol.CLOUDEVENTS, validEvents, invalidEvents, handlerParams.Meta)
			c.JSON(http.StatusOK, response.Ok)
		} else {
			c.JSON(http.StatusBadRequest, response.InvalidContentType)
		}
	}
	return gin.HandlerFunc(fn)
}
