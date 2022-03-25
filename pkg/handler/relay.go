package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

// RelayHandler processes incoming envelopes, splits them in half,
// and sends them to the configured sink. It relies on upstream validation.
func RelayHandler(h EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := envelope.BuildRelayEnvelopesFromRequest(c)
		h.Manifold.Enqueue(envelopes)
	}
	return gin.HandlerFunc(fn)
}
