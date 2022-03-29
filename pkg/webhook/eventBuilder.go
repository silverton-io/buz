package webhook

import (
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const DEFAULT_WEBHOOK_ID string = "f324467b-1213-46a1-b5c8-f173aa5d82e4"

func BuildEvent(c *gin.Context, payload gjson.Result) (WebhookEvent, error) {
	webhookId := c.Param(WEBHOOK_ID_PARAM)
	if webhookId == "" || webhookId == "/" {
		webhookId = DEFAULT_WEBHOOK_ID
	} else {
		webhookId = webhookId[1:]
	}
	event := WebhookEvent{
		Id:      webhookId,
		Payload: payload.Value().(map[string]interface{}),
	}
	return event, nil
}
