package pixel

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/event"
)

const (
	B64_ENCODED_PAYLOAD_PARAM = "hbp"
	ARBITRARY_PIXEL_SCHEMA    = "io.silverton/honeypot/internal/event/pixel/arbitrary/v1.0.json"
)

func BuildEvent(c *gin.Context, params map[string]interface{}) (event.SelfDescribingPayload, error) {
	base64EncodedPayload := params[B64_ENCODED_PAYLOAD_PARAM]
	if base64EncodedPayload != nil {
		p, err := b64.RawStdEncoding.DecodeString(base64EncodedPayload.(string))
		if err != nil {
			return event.SelfDescribingPayload{}, err
		}
		var payload map[string]interface{}
		err = json.Unmarshal(p, &payload)
		if err != nil {
			return event.SelfDescribingPayload{}, err
		}
		e := event.SelfDescribingPayload{
			Schema: ARBITRARY_PIXEL_SCHEMA,
			Data:   payload,
		}
		return e, nil
	}
	e := event.SelfDescribingPayload{
		Schema: ARBITRARY_PIXEL_SCHEMA,
		Data:   params,
	}
	return e, nil
}
