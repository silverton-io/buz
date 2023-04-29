// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package validator

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/protocol"
	"github.com/silverton-io/buz/pkg/registry"
)

// Validate an envelope's payload according to the corresponding schema
func Validate(e envelope.Envelope, registry *registry.Registry) (isValid bool, validationError envelope.ValidationError, schema []byte) {
	// If payload doesn't have a schema associated with it, consider the payload invalid
	if e.Schema == constants.UNKNOWN {
		validationError := envelope.ValidationError{
			ErrorType:       &NoSchemaAssociated.Type,
			ErrorResolution: &NoSchemaAssociated.Resolution,
			Errors:          nil,
		}
		return false, validationError, nil
	}
	schemaExists, schemaContents := registry.Get(e.Schema)
	// If the payload has a schema associated but the schema does not exist
	// in the registry, consider the payload invalid.
	if !schemaExists {
		validationError := envelope.ValidationError{
			ErrorType:       &NoSchemaInBackend.Type,
			ErrorResolution: &NoSchemaInBackend.Resolution,
			Errors:          nil,
		}
		return false, validationError, nil
	} else {
		schema, err := registry.Compiler.Compile(e.Schema)
		if err != nil {
			log.Error().Err(err).Msg("could not compile schema")
			validationError := envelope.ValidationError{
				ErrorType:       &InvalidSchema.Type,
				ErrorResolution: &InvalidSchema.Resolution,
				Errors:          nil,
			}
			return false, validationError, schemaContents
		}
		var payloadToValidate interface{}
		// Snowplow events have to be handled separately, as `self_describing_event` is
		// the only portion that is validated according to a jsonschema.
		if e.Protocol == protocol.SNOWPLOW {
			payloadToValidate = e.Payload["self_describing_event"].(map[string]interface{})["data"]
		} else {
			contents, _ := e.Payload.AsByte()
			err := json.Unmarshal(contents, &payloadToValidate)
			if err != nil {
				log.Error().Err(err).Msg("could not unmarshal payload")
				// If the payload cannot be unmarshaled it should be considered invalid
				validationError := envelope.ValidationError{
					ErrorType:       &InvalidPayload.Type,
					ErrorResolution: &InvalidPayload.Resolution,
					Errors:          nil,
				}
				return false, validationError, nil
			}
		}
		// If the payload is not present at all it should be considered invalid.
		if payloadToValidate == nil {
			validationError := envelope.ValidationError{
				ErrorType:       &PayloadNotPresent.Type,
				ErrorResolution: &PayloadNotPresent.Resolution,
				Errors:          nil,
			}
			return false, validationError, nil
		}
		startTime := time.Now().UTC()
		vErr := schema.Validate(payloadToValidate)
		log.Debug().Msg("🟡 event validated in " + time.Now().UTC().Sub(startTime).String())
		jsonschemaValidationErr := &jsonschema.ValidationError{}
		if vErr != nil {
			if errors.As(vErr, &jsonschemaValidationErr) {
				var validationErrs = []envelope.PayloadValidationError{}
				for _, cause := range jsonschemaValidationErr.Causes {
					validationErr := envelope.PayloadValidationError{
						Field:       cause.InstanceLocation,
						Description: cause.Message,
						ErrorType:   cause.KeywordLocation,
					}
					validationErrs = append(validationErrs, validationErr)
				}
				validationError := envelope.ValidationError{
					ErrorType:       &InvalidPayload.Type,
					ErrorResolution: &InvalidPayload.Resolution,
					Errors:          validationErrs,
				}
				return false, validationError, schemaContents
			}
			return false, envelope.ValidationError{}, schemaContents
		}
		return true, envelope.ValidationError{}, schemaContents
	}
}
