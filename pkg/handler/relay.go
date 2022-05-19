package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/params"
	"github.com/silverton-io/honeypot/pkg/response"
)

// RelayHandler processes incoming envelopes, splits them in half,
// and sends them to the configured sink. It relies on upstream validation.
func RelayHandler(h params.Handler) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := envelope.BuildRelayEnvelopesFromRequest(c, h.CollectorMeta)
		err := h.Manifold.Distribute(envelopes, *h.ProtocolStats)
		if err != nil {
			c.Header("Retry-After", response.RETRY_AFTER_3)
			c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
		} else {
			c.JSON(200, response.Ok)
		}
	}
	return gin.HandlerFunc(fn)
}
