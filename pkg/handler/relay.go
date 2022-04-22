package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/response"
)

// RelayHandler processes incoming envelopes, splits them in half,
// and sends them to the configured sink. It relies on upstream validation.
func RelayHandler(h EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := envelope.BuildRelayEnvelopesFromRequest(c)
		err := h.Manifold.Distribute(envelopes)
		if err != nil {
			c.Request.Header.Add("Retry-After", response.RETRY_AFTER_60)
			c.JSON(http.StatusServiceUnavailable, response.DistributionError)
		} else {
			c.JSON(200, response.Ok)
		}
	}
	return gin.HandlerFunc(fn)
}
