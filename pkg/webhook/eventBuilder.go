package webhook

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/tidwall/gjson"
)

const ARBITRARY_WEBHOOK_SCHEMA = "io.silverton/honeypot/internal/event/webhook/arbitrary/v1.0.json"

func BuildEvent(c *gin.Context, payload gjson.Result) (event.SelfDescribingPayload, error) {
	schemaName := util.GetSchemaNameFromRequest(c, ARBITRARY_WEBHOOK_SCHEMA)
	e := event.SelfDescribingPayload{
		Schema: schemaName,
		Data:   payload.Value().(map[string]interface{}),
	}
	return e, nil
}
