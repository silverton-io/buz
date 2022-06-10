package webhook

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/tidwall/gjson"
)

const ARBITRARY_WEBHOOK_SCHEMA = "io.silverton/honeypot/internal/event/webhook/arbitrary/v1.0.json"

func BuildEvent(c *gin.Context, payload gjson.Result) (event.SelfDescribingPayload, error) {
	e := event.SelfDescribingPayload{
		Schema: ARBITRARY_WEBHOOK_SCHEMA,
		Data:   payload.Value().(map[string]interface{}),
	}
	return e, nil
}
