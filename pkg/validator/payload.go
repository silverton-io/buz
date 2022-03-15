package validator

import (
	"context"
	"encoding/json"
	"time"

	"github.com/qri-io/jsonschema"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/event"
)

func validatePayload(payload []byte, schema []byte) (isValid bool, validationError event.ValidationError) {
	ctx := context.Background()
	startTime := time.Now()
	s := &jsonschema.Schema{}
	unmarshalErr := json.Unmarshal(schema, s)
	if unmarshalErr != nil {
		log.Error().Stack().Err(unmarshalErr).Msg("failed to unmarshal schema")
	}
	validationErrs, vErr := s.ValidateBytes(ctx, payload)

	if unmarshalErr != nil || vErr != nil {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		validationError := event.ValidationError{
			ErrorType:       "invalid schema",
			ErrorResolution: "ensure schema is properly formatted",
			Errors:          nil,
		}
		return false, validationError
	}
	if len(validationErrs) == 0 {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return true, event.ValidationError{}
	} else {
		var payloadValidationErrors []event.PayloadValidationError
		for _, validationErr := range validationErrs {
			payloadValidationError := event.PayloadValidationError{
				Field:       validationErr.PropertyPath,
				Description: validationErr.Message,
				ErrorType:   validationErr.Error(),
			}
			payloadValidationErrors = append(payloadValidationErrors, payloadValidationError)
		}
		validationError := event.ValidationError{
			ErrorType:       "invalid payload",
			ErrorResolution: "correct payload format",
			Errors:          payloadValidationErrors,
		}
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return false, validationError
	}
}
