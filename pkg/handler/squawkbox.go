package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func SquawkboxHandler(h EventHandlerParams, eventProtocol string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var envelopes []envelope.Envelope
		switch eventProtocol {
		case protocol.SNOWPLOW:
			envelopes = envelope.BuildSnowplowEnvelopesFromRequest(c, *h.Config)
		case protocol.CLOUDEVENTS:
			envelopes = envelope.BuildCloudeventEnvelopesFromRequest(c, *h.Config)
		case protocol.GENERIC:
			envelopes = envelope.BuildGenericEnvelopesFromRequest(c, *h.Config)
		}
		annotatedEnvelopes := validator.Annotate(envelopes, h.Cache)
		c.JSON(http.StatusOK, annotatedEnvelopes)
	}
	return gin.HandlerFunc(fn)
}
