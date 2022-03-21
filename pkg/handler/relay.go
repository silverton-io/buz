package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/validator"
)

// RelayHandler processes incoming envelopes, splits them in half,
// and sends them to the configured sink. It relies on upstream validation.
func RelayHandler(p EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := envelope.BuildRelayEnvelopesFromRequest(c)
		validEnvelopes, invalidEnvelopes := validator.Bifurcate(envelopes)
		p.Sink.BatchPublishValidAndInvalid(ctx, validEnvelopes, invalidEnvelopes, p.Meta)
	}
	return gin.HandlerFunc(fn)
}
