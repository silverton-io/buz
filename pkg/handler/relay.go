package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/validator"
)

// RelayHandler processes incoming envelopes, splits them in half,
// and sends them to the configured sink. It relies on upstream validation.
func RelayHandler(h EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := envelope.BuildRelayEnvelopesFromRequest(c)
		annotatedEnvelopes := validator.Annotate(envelopes, h.Cache)
		h.Manifold.Enqueue(annotatedEnvelopes)
	}
	return gin.HandlerFunc(fn)
}
