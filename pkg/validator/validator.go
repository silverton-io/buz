package validator

import (
	"context"
	"encoding/json"
	"time"

	"github.com/qri-io/jsonschema"
	"github.com/rs/zerolog/log"
)

type PayloadValidationError struct {
	Field       string `json:"field"`
	Description string `json:"description"`
	ErrorType   string `json:"errorType"`
}

type ValidationError struct {
	ErrorType       string                    `json:"errorType"`
	ErrorResolution string                    `json:"errorResolution"`
	Errors          *[]PayloadValidationError `json:"payloadValidationErrors"`
}

func ValidatePayload(payload map[string]interface{}, schema []byte) (isValid bool, validationError ValidationError) {
	ctx := context.Background()
	startTime := time.Now()
	s := &jsonschema.Schema{}
	json.Unmarshal(schema, s)
	data, _ := json.Marshal(payload)
	validationErrs, err := s.ValidateBytes(ctx, data)

	if err != nil {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		validationError := ValidationError{
			ErrorType:       "invalid schema",
			ErrorResolution: "ensure schema is properly-formatted",
			Errors:          nil,
		}
		return false, validationError
	}
	if len(validationErrs) == 0 {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return true, ValidationError{}
	} else {
		var payloadValidationErrors []PayloadValidationError
		for _, validationErr := range validationErrs {
			payloadValidationError := PayloadValidationError{
				Field:       validationErr.PropertyPath,
				Description: validationErr.Message,
				ErrorType:   validationErr.Error(),
			}
			payloadValidationErrors = append(payloadValidationErrors, payloadValidationError)
		}
		validationError := ValidationError{
			ErrorType:       "invalid payload",
			ErrorResolution: "correct payload format",
			Errors:          &payloadValidationErrors,
		}
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return false, validationError
	}
}
