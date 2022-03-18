package webhook

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

const DEFAULT_WEBHOOK_ID string = "f324467b-1213-46a1-b5c8-f173aa5d82e4"

func BuildEvent(c *gin.Context, payload gjson.Result) (WebhookEvent, error) {
	webhookId := c.Param(WEBHOOK_ID_PARAM)[1:]
	if webhookId == "" {
		webhookId = DEFAULT_WEBHOOK_ID
	}
	event := WebhookEvent{}
	event.Id = webhookId
	pBytes, err := json.Marshal(payload.Value().(map[string]interface{}))
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not marshal webhook payload")
		return WebhookEvent{}, nil
	}
	err = json.Unmarshal(pBytes, &event)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not unmarshal payload to WebhookEvent")
		return WebhookEvent{}, nil
	}
	return event, nil
}
