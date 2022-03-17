package cloudevents

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

func BuildEvent(payload gjson.Result) (CloudEvent, error) {
	event := CloudEvent{}
	event.DataContentType = "application/json" // Always

	pBytes, err := json.Marshal(payload.Value().(map[string]interface{}))
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not marshal cloudevent payload")
		return CloudEvent{}, nil
	}
	err = json.Unmarshal(pBytes, &event)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not unmarshal payload to CloudEvent")
		return CloudEvent{}, err
	}
	return event, nil
}
