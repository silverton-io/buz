package validator

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/silverton-io/honeypot/pkg/cache"
)

func validateEvent(event event.Event, cache *cache.SchemaCache) (isValid bool, validationError ValidationError) {
	schemaName := event.Context.GetDataSchema()
	if schemaName == "" {
		validationError := ValidationError{ // Enforce dataschema is present in all cloudevents
			ErrorType:       "cloudevent missing dataschema",
			ErrorResolution: "publish the cloudevent with dataschema context",
			Errors:          nil,
		}
		return false, validationError
	}
	schemaExists, schemaContents := cache.Get(schemaName)
	if !schemaExists {
		validationError := ValidationError{
			ErrorType:       "nonexistent schema",
			ErrorResolution: "publish the specified schema to the cache backend",
			Errors:          nil,
		}
		return false, validationError
	} else {
		isValid, validationError := ValidatePayload(event.Data(), schemaContents)
		return isValid, validationError
	}
}
