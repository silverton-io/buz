// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package validator

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/protocol"
	"github.com/silverton-io/buz/pkg/registry"
	"github.com/silverton-io/buz/pkg/util"
)

func ValidatePayload(e envelope.Envelope, registry *registry.Registry) (isValid bool, validationError envelope.ValidationError, schema []byte) {
	// FIXME- Short-circuit if the event is an unknown event
	if e.Schema == "" {
		validationError := envelope.ValidationError{
			ErrorType:       &NoSchemaAssociated.Type,
			ErrorResolution: &NoSchemaAssociated.Resolution,
			Errors:          nil,
		}
		return false, validationError, nil
	}
	schemaExists, schemaContents := registry.Get(e.Schema)
	if !schemaExists {
		validationError := envelope.ValidationError{
			ErrorType:       &NoSchemaInBackend.Type,
			ErrorResolution: &NoSchemaInBackend.Resolution,
			Errors:          nil,
		}
		return false, validationError, nil
	} else {
		var payloadToValidate []byte
		if e.Protocol != protocol.SNOWPLOW {
			payloadToValidate, _ = e.Payload.AsByte()
		} else { // If the envelope protocol is snowplow, only validate SelfDescribingEvent
			e := e.Payload["self_describing_event"].(map[string]interface{})["data"]
			util.Pprint(e)
			payloadToValidate, _ = json.Marshal(e)
		}

		payload, err := e.Payload.AsByte()
		if err != nil {
			log.Error().Stack().Err(err).Msg("🔴 could not marshal payload")
			validationError := envelope.ValidationError{
				ErrorType:       &InvalidPayload.Type,
				ErrorResolution: &InvalidPayload.Resolution,
				Errors:          nil,
			}
			return false, validationError, nil
		}
		if payload == nil {
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
