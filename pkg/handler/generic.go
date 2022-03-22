package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func GenericHandler(h EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := envelope.BuildGenericEnvelopesFromRequest(c, *h.Config)
		validEvents, invalidEvents, stats := validator.BifurcateAndAnnotate(envelopes, h.Cache)
		h.Sink.BatchPublishValidAndInvalid(ctx, validEvents, invalidEvents)
		util.Pprint("published")
		h.Meta.ProtocolStats.Merge(&stats)
		util.Pprint("merged stats")
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
