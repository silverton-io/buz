package pixel

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/util"
)

const (
	B64_ENCODED_PAYLOAD_PARAM string = "hbp"
	ARBITRARY_PIXEL_SCHEMA    string = "io.silverton/honeypot/internal/event/pixel/arbitrary/v1.0.json"
)

func BuildEvent(c *gin.Context) (event.SelfDescribingPayload, error) {
	params := util.MapUrlParams(c)
	schemaName := util.GetSchemaNameFromRequest(c, ARBITRARY_PIXEL_SCHEMA)
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
			Schema: *schemaName,
			Data:   payload,
		}
		return e, nil
	}
	e := event.SelfDescribingPayload{
		Schema: *schemaName,
		Data:   params,
	}
	return e, nil
}
