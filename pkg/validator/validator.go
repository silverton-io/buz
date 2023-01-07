// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package validator

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/event"
)

func ValidatePayload(schemaContents []byte, payLoad event.Payload) (isValid bool, validationError envelope.ValidationError) {
	payload, err := payLoad.AsByte()
	if err != nil {
		log.Error().Stack().Err(err).Msg("🔴 could not marshal payload")
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
