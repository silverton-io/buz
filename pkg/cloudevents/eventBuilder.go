// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

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
