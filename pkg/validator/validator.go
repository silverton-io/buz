// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package validator

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/registry"
)

func ValidatePayload(schemaName string, payload envelope.Payload, registry *registry.Registry) (isValid bool, validationError envelope.ValidationError, schema []byte) {
	// FIXME- Short-circuit if the event is an unknown event
	if schemaName == "" {
		validationError := envelope.ValidationError{
			ErrorType:       &InvalidPayload.Type,
			ErrorResolution: &InvalidPayload.Resolution,
			Errors:          nil,
		}
		return false, validationError
	}
	if payload == nil {
		validationError := envelope.ValidationError{
			ErrorType:       &PayloadNotPresent.Type,
			ErrorResolution: &PayloadNotPresent.Resolution,
			Errors:          nil,
		}
		return false, validationError
	}
	isValid, validationError = validatePayload(payload, schemaContents)
	return isValid, validationError

}
