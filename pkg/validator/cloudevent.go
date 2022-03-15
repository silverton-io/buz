package validator

import (
	ce "github.com/cloudevents/sdk-go/v2/event"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/event"
)

func validateCloudEvent(cloudevent ce.Event, cache *cache.SchemaCache) (isValid bool, validationError event.ValidationError, schema []byte) {
	schemaName := cloudevent.Context.GetDataSchema()
	if schemaName == "" {
		validationError := event.ValidationError{ // Enforce dataschema is present in all cloudevents
			ErrorType:       "cloudevent missing dataschema",
			ErrorResolution: "publish the cloudevent with dataschema context",
			Errors:          nil,
		}
		return false, validationError, nil
	}
	schemaExists, schemaContents := cache.Get(schemaName)
	if !schemaExists {
		validationError := event.ValidationError{
			ErrorType:       "nonexistent cloudevent schema",
			ErrorResolution: "publish the specified schema to the cache backend",
			Errors:          nil,
		}
		return false, validationError, nil
	} else {
		isValid, validationError := validatePayload(cloudevent.Data(), schemaContents)
		return isValid, validationError, schemaContents
	}
}
