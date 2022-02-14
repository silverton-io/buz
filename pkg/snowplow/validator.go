package snowplow

import (
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/validator"
)

const (
	IGLU = "iglu"
)

func ValidateEvent(event Event, cache *cache.SchemaCache) (isValid bool, validationError validator.ValidationError, schema []byte) {
	if event.Event == SELF_DESCRIBING_EVENT { // Only validate self describing events
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
	} else {
		return true, validator.ValidationError{}, nil
	}
}
