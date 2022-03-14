package ce

import (
	ce "github.com/cloudevents/sdk-go/v2/event"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func validateEvent(event ce.Event, cache *cache.SchemaCache) (isValid bool, validationError validator.ValidationError) {
	schemaName := event.Context.GetDataSchema()
	if schemaName == "" {
		validationError := validator.ValidationError{ // Enforce dataschema is present in all cloudevents
			ErrorType:       "cloudevent missing dataschema",
			ErrorResolution: "publish the cloudevent with dataschema context",
			Errors:          nil,
		}
		return false, validationError
	}
	schemaExists, schemaContents := cache.Get(schemaName)
	if !schemaExists {
		validationError := validator.ValidationError{
			ErrorType:       "nonexistent schema",
			ErrorResolution: "publish the specified schema to the cache backend",
			Errors:          nil,
		}
		return false, validationError
	} else {
		isValid, validationError := validator.ValidatePayload(event.Data(), schemaContents)
		return isValid, validationError
	}
}
