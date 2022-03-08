package validator

import (
	"bytes"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

type PayloadValidationError struct {
	Field       string `json:"field"`
	Context     string `json:"context"`
	Description string `json:"description"`
	ErrorType   string `json:"errorType"`
}

type ValidationError struct {
	ErrorType       string                    `json:"errorType"`
	ErrorResolution string                    `json:"errorResolution"`
	Errors          *[]PayloadValidationError `json:"payloadValidationErrors"`
}

func ValidatePayload(payload map[string]interface{}, schema []byte) (isValid bool, validationError ValidationError) {
	startTime := time.Now()

	compiler := jsonschema.NewCompiler()
	err := compiler.AddResource("schema.json", bytes.NewBuffer(schema))
	s, _ := compiler.Compile("schema.json")
	validationErr := s.Validate(payload)

	docLoader := gojsonschema.NewGoLoader(payload)
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		validationError := ValidationError{
			ErrorType:       "invalid schema",
			ErrorResolution: "ensure schema is properly-formatted",
			Errors:          nil,
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
		validationError := ValidationError{
			ErrorType:       "invalid payload",
			ErrorResolution: "correct payload format",
			Errors:          &payloadValidationErrors,
		}
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return false, validationError
	}
}
