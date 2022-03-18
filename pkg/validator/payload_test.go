package validator

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qri-io/jsonschema"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type output struct {
	isValid         bool
	validationError envelope.ValidationError
}

func generatePayloadValidationErrs(payload []byte, schema []byte) []envelope.PayloadValidationError {
	ctx := context.Background()
	s := &jsonschema.Schema{}
	json.Unmarshal(schema, s)
	validationErrs, _ := s.ValidateBytes(ctx, payload)

	var payloadValidationErrors []envelope.PayloadValidationError
	for _, validationErr := range validationErrs {
		payloadValidationError := envelope.PayloadValidationError{
			Field:       validationErr.PropertyPath,
			Description: validationErr.Message,
			ErrorType:   validationErr.Error(),
		}
		payloadValidationErrors = append(payloadValidationErrors, payloadValidationError)
	}
	return payloadValidationErrors
}

func TestValidatePayload(t *testing.T) {

	validPayload := []byte(`{"id": 10, "action": "did"}`)
	invalidPayload := []byte(`{"id": 10, "action": "did", "somethingBad": 10}`)
	validSchema := []byte(`
	{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$id": "did",
		"title": "did",
		"type": "object",
		"properties": {
			"id": {
				"type": "number"
			},
			"action": {
				"type": "string"
			}
		},
		"required": ["id", "action"],
		"additionalProperties": false
	}
	`)
	invalidSchema := []byte(`{"something": yup`)

	invalidPayloadValidationErrs := generatePayloadValidationErrs(invalidPayload, validSchema)

	var testCases = []struct {
		name    string
		payload []byte
		schema  []byte
		want    output
	}{
		{"valid payload valid schema", validPayload, validSchema, output{true, envelope.ValidationError{}}},
		{"valid payload invalid schema", validPayload, invalidSchema, output{false, envelope.ValidationError{ErrorType: &InvalidSchema.Type, ErrorResolution: &InvalidSchema.Resolution, Errors: nil}}},
		{"invalid payload valid schema", invalidPayload, validSchema, output{false, envelope.ValidationError{ErrorType: &InvalidPayload.Type, ErrorResolution: &InvalidPayload.Resolution, Errors: invalidPayloadValidationErrs}}},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			isValid, vErr := validatePayload(tc.payload, tc.schema)
			if isValid != tc.want.isValid {
				t.Fatalf(`got %v, want %v`, isValid, tc.want.isValid)
			}
			errEquiv := reflect.DeepEqual(vErr, tc.want.validationError)
			if !errEquiv {
				t.Fatalf(`got %v, want %v`, vErr, tc.want.validationError)
			}
		})
	}
}
