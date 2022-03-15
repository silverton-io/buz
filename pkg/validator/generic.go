package validator

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/tidwall/gjson"
)

func validateGenericEvent(genericEvent gjson.Result, schemaName string, cache *cache.SchemaCache, conf *config.Generic) (isValid bool, validationError event.ValidationError, schema []byte) {
	if genericEvent.Value() == nil {
		validationError := event.ValidationError{
			ErrorType:       "payload not present at " + conf.Payload.RootKey + "." + conf.Payload.DataKey,
			ErrorResolution: "generic payload configuration and payload path should match",
			Errors:          nil,
		}
		return false, validationError, nil
	}
	if schemaName == "" { // Event does not have schema associated - always invalid.
		validationError := event.ValidationError{
			ErrorType:       "schema not present at " + conf.Payload.RootKey + "." + conf.Payload.SchemaKey,
			ErrorResolution: "generic schema configuration and schema path should match",
			Errors:          nil,
		}
		return false, validationError, nil
	}
	schemaExists, schemaContents := cache.Get(schemaName)
	if !schemaExists { // Referenced schema is not present in the cache backend - always invalid
		validationError := event.ValidationError{
			ErrorType:       "nonexistent schema",
			ErrorResolution: "publish the specified schema to the cache backend",
			Errors:          nil,
		}
		return false, validationError, nil
	} else {
		isValid, validationError := validatePayload([]byte(genericEvent.Raw), schemaContents)
		return isValid, validationError, schemaContents
	}
}
