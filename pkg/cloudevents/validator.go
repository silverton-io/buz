package ce

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/validator"
	"github.com/tidwall/gjson"
)

func validateEvent(event event.Event, cache *cache.SchemaCache) (isValid bool, validationError validator.ValidationError) {
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
		payload := gjson.ParseBytes(event.Data()).Value().(map[string]interface{})
		isValid, validationError := validator.ValidatePayload(payload, schemaContents)
		return isValid, validationError
	}
}
