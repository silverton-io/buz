package snowplow

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/validator"
)

const (
	IGLU = "iglu"
)

func validateEvent(event Event, cache *cache.SchemaCache) (isValid bool, validationError validator.ValidationError, schema []byte) {
	switch event.Event {
	case UNKNOWN_EVENT:
		validationError := validator.ValidationError{
			ErrorType:       "unknown event type",
			ErrorResolution: "event type needs to adhere to the snowplow tracker protocol",
			Errors:          nil,
		}
		return false, validationError, nil
	case SELF_DESCRIBING_EVENT:
		schemaName := event.Self_describing_event.Schema
		if schemaName[:4] == IGLU { // If schema path starts with iglu:, get rid of it.
			schemaName = schemaName[5:]
		}
		schemaExists, schemaContents := cache.Get(schemaName)
		if !schemaExists {
			validationError := validator.ValidationError{
				ErrorType:       "nonexistent schema",
				ErrorResolution: "publish the specified schema to the cache backend",
				Errors:          nil,
			}
			return false, validationError, nil
		} else {
			eventData, err := json.Marshal(event.Self_describing_event.Data)
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not marshal event")
			}
			isValid, validationError := validator.ValidatePayload(eventData, schemaContents)
			return isValid, validationError, schemaContents
		}
	default:
		return true, validator.ValidationError{}, nil // Treat non-self-describing events as "valid"
	}
}
