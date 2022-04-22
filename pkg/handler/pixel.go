package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/annotator"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/response"
)

const PX string = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8Xw8AAoMBgDTD2qgAAAAASUVORK5CYII="

func PixelHandler(h EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := envelope.BuildPixelEnvelopesFromRequest(c)
		annotatedEnvelopes := annotator.Annotate(envelopes, h.Cache)
		err := h.Manifold.Distribute(annotatedEnvelopes, h.Meta)
		if err != nil {
			c.Header("Retry-After", response.RETRY_AFTER_60)
			c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
		} else {
			b, _ := base64.StdEncoding.DecodeString(PX)
			c.Data(http.StatusOK, "image/png", b)
		}
	}
	return gin.HandlerFunc(fn)
}
