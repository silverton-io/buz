package webhook

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

func BuildEvent(payload gjson.Result) (WebhookEvent, error) {
	event := WebhookEvent{}
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
