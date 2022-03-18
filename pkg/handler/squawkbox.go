package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func SquawkboxHandler(p EventHandlerParams, eventProtocol string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var envelopes []envelope.Envelope
		switch eventProtocol {
		case protocol.SNOWPLOW:
			envelopes = envelope.BuildSnowplowEnvelopesFromRequest(c, *p.Config)
		case protocol.CLOUDEVENTS:
			envelopes = envelope.BuildCloudeventEnvelopesFromRequest(c, *p.Config)
		case protocol.GENERIC:
			envelopes = envelope.BuildGenericEnvelopesFromRequest(c, *p.Config)
		}
		validEnvelopes, invalidEnvelopes := validator.BifurcateAndAnnotate(envelopes, p.Cache)
		envelopes = append(validEnvelopes, invalidEnvelopes...)
		c.JSON(http.StatusOK, envelopes)
	}
	return gin.HandlerFunc(fn)
}
