// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/annotator"
	"github.com/silverton-io/buz/pkg/envelope"
	cloudevents "github.com/silverton-io/buz/pkg/inputCloudevents"
	pixel "github.com/silverton-io/buz/pkg/inputPixel"
	selfdescribing "github.com/silverton-io/buz/pkg/inputSelfDescribing"
	webhook "github.com/silverton-io/buz/pkg/inputWebhook"
	"github.com/silverton-io/buz/pkg/params"
	"github.com/silverton-io/buz/pkg/protocol"
)

func SquawkboxHandler(h params.Handler, eventProtocol string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var envelopes []envelope.Envelope
		switch eventProtocol {
		case protocol.SNOWPLOW:
			envelopes = envelope.BuildSnowplowEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		case protocol.CLOUDEVENTS:
			envelopes = cloudevents.BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		case protocol.SELF_DESCRIBING:
			envelopes = selfdescribing.BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		case protocol.PIXEL:
			envelopes = pixel.BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		case protocol.WEBHOOK:
			envelopes = webhook.BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		}
		annotatedEnvelopes := annotator.Annotate(envelopes, h.Registry)
		c.JSON(http.StatusOK, annotatedEnvelopes)
	}
	return gin.HandlerFunc(fn)
}
