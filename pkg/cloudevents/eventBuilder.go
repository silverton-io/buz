package cloudevents

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func BuildEvent(c *gin.Context, payload map[string]interface{}) CloudEvent {
	event := CloudEvent{}
	pBytes, err := json.Marshal(payload)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not marshal cloudevent payload")
	}
	err = json.Unmarshal(pBytes, &event)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not unmarshal payload to CloudEvent")
	}
	return event
}
