package pixel

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

const UNNAMED_PIXEL_ID = "unnamed"
const B64_ENCODED_PAYLOAD_PARAM = "hbp"

func BuildEvent(c *gin.Context, params map[string]interface{}) (PixelEvent, error) {
	base64EncodedPayload := params[B64_ENCODED_PAYLOAD_PARAM]
	if base64EncodedPayload != nil {
		p, err := b64.RawStdEncoding.DecodeString(base64EncodedPayload.(string))
		if err != nil {
			return PixelEvent{}, err
		}
		var payload map[string]interface{}
		err = json.Unmarshal(p, &payload)
		if err != nil {
			return PixelEvent{}, err
		}
		event := PixelEvent{
			Id:      UNNAMED_PIXEL_ID, // FIXME!
			Payload: payload,
		}
		return event, nil
	}
	event := PixelEvent{
		Id:      UNNAMED_PIXEL_ID, // FIXME!
		Payload: params,
	}
	return event, nil
}
