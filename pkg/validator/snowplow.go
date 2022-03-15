package validator

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/snowplow"
)

func validateSnowplowEvent(spEvent snowplow.SnowplowEvent, cache *cache.SchemaCache) (isValid bool, validationError event.ValidationError, schema []byte) {
	switch spEvent.Event {
	case snowplow.UNKNOWN_EVENT:
		validationError := event.ValidationError{
			ErrorType:       "unknown event type",
			ErrorResolution: "event type needs to adhere to the snowplow tracker protocol",
			Errors:          nil,
		}
		return false, validationError, nil
	case snowplow.SELF_DESCRIBING_EVENT:
		schemaName := spEvent.Self_describing_event.Schema
		if schemaName[:4] == snowplow.IGLU { // If schema path starts with iglu:, get rid of it.
			schemaName = schemaName[5:]
		}
		schemaExists, schemaContents := cache.Get(schemaName)
		if !schemaExists {
			validationError := event.ValidationError{
				ErrorType:       "nonexistent schema",
				ErrorResolution: "publish the specified schema to the cache backend",
				Errors:          nil,
			}
			return false, validationError, nil
		} else {
			eventData, err := json.Marshal(spEvent.Self_describing_event.Data)
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not marshal event")
			}
			isValid, validationError := validatePayload(eventData, schemaContents)
			return isValid, validationError, schemaContents
		}
	default:
		return true, event.ValidationError{}, nil // Treat non-self-describing events as "valid"
	}
}
