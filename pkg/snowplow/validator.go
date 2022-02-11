package snowplow

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/xeipuuv/gojsonschema"
)

const (
	IGLU = "iglu"
)

func validateSelfDescribingPayload(payload SelfDescribingPayload, cache *cache.SchemaCache) (isValid bool, validationError ValidationError, schema []byte) {
	log.Debug().Msg("validating self describing payload")
	startTime := time.Now()
	schemaName := payload.Schema
	if schemaName[:4] == IGLU { // If schema path starts with iglu, get rid of it.
		schemaName = schemaName[5:]
	}
	fmt.Println(schemaName)
	schemaExists, schema := cache.Get(schemaName)
	// Short-circuit if schema can't be found in either cache or remote backend.
	if !schemaExists {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		errorType := "nonexistent schema"
		validationError := ValidationError{
			ErrorType: &errorType,
			Errors:    nil,
		}
		return false, validationError, nil
	}
	docLoader := gojsonschema.NewGoLoader(payload.Data)
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not validate payload " + payload.Schema) // TODO! Test an actual invalid schema here
		errorType := "invalid schema format"
		validationError := ValidationError{
			ErrorType: &errorType,
			Errors:    nil,
		}
		return false, validationError, nil
	}
	if result.Valid() {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return true, ValidationError{}, schema
	} else {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		var payloadValidationErrors []PayloadValidationError

		for _, validationError := range result.Errors() {
			payloadValidationError := PayloadValidationError{
				Field:       validationError.Field(),
				Context:     validationError.Context().String(),
				Description: validationError.Description(),
				ErrorType:   validationError.Type(),
			}
			payloadValidationErrors = append(payloadValidationErrors, payloadValidationError)
		}
		errorType := "invalid payload"
		validationError := ValidationError{
			ErrorType: &errorType,
			Errors:    &payloadValidationErrors,
		}
		return false, validationError, schema
	}
}

func ValidateEvent(event Event, cache *cache.SchemaCache) (isValid bool, validationError ValidationError, schema []byte) {
	if event.Event == SELF_DESCRIBING_EVENT { // Only validate if event type is self describing
		return validateSelfDescribingPayload(SelfDescribingPayload(*event.Self_describing_event), cache)
	} else {
		return true, ValidationError{}, nil
	}
}
