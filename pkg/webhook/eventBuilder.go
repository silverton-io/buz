package webhook

import (
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const UNNAMED_WEBHOOK_ID string = "unnamed"

func BuildEvent(c *gin.Context, payload gjson.Result) (WebhookEvent, error) {
	webhookId := c.Param(WEBHOOK_ID_PARAM)
	if webhookId == "" || webhookId == "/" {
		webhookId = UNNAMED_WEBHOOK_ID
	} else {
		webhookId = webhookId[1:]
	}
	event := WebhookEvent{
		Id:      webhookId,
		Payload: payload.Value().(map[string]interface{}),
	}
	return event, nil
}
