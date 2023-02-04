// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package handler

// func SquawkboxHandler(h params.Handler, eventProtocol string) gin.HandlerFunc {
// 	fn := func(c *gin.Context) {
// 		var envelopes []envelope.Envelope
// 		switch eventProtocol {
// 		case protocol.SNOWPLOW:
// 			envelopes = snowplow.BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
// 		case protocol.CLOUDEVENTS:
// 			envelopes = cloudevents.BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
// 		case protocol.SELF_DESCRIBING:
// 			envelopes = selfdescribing.BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
// 		case protocol.PIXEL:
// 			envelopes = pixel.BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
// 		case protocol.WEBHOOK:
// 			envelopes = webhook.BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
// 		}
// 		annotatedEnvelopes := annotator.Annotate(envelopes, h.Registry)
// 		c.JSON(http.StatusOK, annotatedEnvelopes)
// 	}
// 	return gin.HandlerFunc(fn)
// }
