// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package validator

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
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
		var payloadToValidate []byte
		var err error
		// Snowplow events have to be handled separately, as `self_describing_event` is
		// the only portion that is validated according to a jsonschema.
		if e.Protocol == protocol.SNOWPLOW {
			e := e.Payload["self_describing_event"].(map[string]interface{})["data"]
			payloadToValidate, err = json.Marshal(e)
		} else {
			payloadToValidate, err = e.Payload.AsByte()
		}
		// If the payload cannot be marshaled it should be considered invalid.
		if err != nil {
			log.Error().Stack().Err(err).Msg("🔴 could not marshal payload")
			validationError := envelope.ValidationError{
				ErrorType:       &InvalidPayload.Type,
				ErrorResolution: &InvalidPayload.Resolution,
				Errors:          nil,
			}
			return false, validationError, nil
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
		isValid, validationError := validatePayload(payloadToValidate, schemaContents)
		return isValid, validationError, schemaContents
	}
}
