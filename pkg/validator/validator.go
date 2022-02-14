package validator

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

type PayloadValidationError struct {
	Field       string `json:"field"`
	Context     string `json:"context"`
	Description string `json:"description"`
	ErrorType   string `json:"errorType"`
}

type ValidationError struct {
	ErrorType *string                   `json:"errorType"`
	Errors    *[]PayloadValidationError `json:"payloadValidationErrors"`
}

func ValidatePayload(payload map[string]interface{}, schema []byte) (isValid bool, validationError ValidationError) {
	startTime := time.Now()
	docLoader := gojsonschema.NewGoLoader(payload)
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		errorType := "invalid schema"
		validationError := ValidationError{
			ErrorType: &errorType,
			Errors:    nil,
		}
		return false, validationError
	}
	if result.Valid() {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return true, ValidationError{}
	} else {
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
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return false, validationError
	}
}
