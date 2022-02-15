package snowplow

import (
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/validator"
)

const (
	IGLU = "iglu"
)

func ValidateEvent(event Event, cache *cache.SchemaCache) (isValid bool, validationError validator.ValidationError, schema []byte) {
	switch event.Event {
	case UNKNOWN_EVENT:
		errorType := "unknown event type"
		validationError := validator.ValidationError{
			ErrorType: &errorType,
			Errors:    nil,
		}
		return false, validationError, nil
	case SELF_DESCRIBING_EVENT:
		schemaName := event.Self_describing_event.Schema
		if schemaName[:4] == IGLU { // If schema path starts with iglu:, get rid of it.
			schemaName = schemaName[5:]
		}
		schemaExists, schemaContents := cache.Get(schemaName)
		if !schemaExists {
			errorType := "nonexistent schema"
			validationError := validator.ValidationError{
				ErrorType: &errorType,
				Errors:    nil,
			}
			return false, validationError, nil
		} else {
			isValid, validationError := validator.ValidatePayload(event.Self_describing_event.Data, schemaContents)
			return isValid, validationError, schemaContents
		}
	default:
		return true, validator.ValidationError{}, nil // Treat non-self-describing events as "valid"
	}
}
