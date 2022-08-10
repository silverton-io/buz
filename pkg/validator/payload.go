// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package validator

import (
	"context"
	"encoding/json"
	"time"

	"github.com/qri-io/jsonschema"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

func validatePayload(payload []byte, schema []byte) (isValid bool, validationError envelope.ValidationError) {
	ctx := context.Background()
	startTime := time.Now().UTC()
	s := &jsonschema.Schema{}
	unmarshalErr := json.Unmarshal(schema, s)
	if unmarshalErr != nil {
		log.Error().Stack().Err(unmarshalErr).Msg("failed to unmarshal schema")
	}
	validationErrs, vErr := s.ValidateBytes(ctx, payload)

	if unmarshalErr != nil || vErr != nil {
		log.Debug().Msg("event validated in " + time.Now().UTC().Sub(startTime).String())
		validationError := envelope.ValidationError{
			ErrorType:       &InvalidSchema.Type,
			ErrorResolution: &InvalidSchema.Resolution,
			Errors:          nil,
		}
		return false, validationError
	}
	if len(validationErrs) == 0 {
		log.Debug().Msg("event validated in " + time.Now().UTC().Sub(startTime).String())
		return true, envelope.ValidationError{}
	} else {
		var payloadValidationErrors []envelope.PayloadValidationError
		for _, validationErr := range validationErrs {
			payloadValidationError := envelope.PayloadValidationError{
				Field:       validationErr.PropertyPath,
				Description: validationErr.Message,
				ErrorType:   validationErr.Error(),
			}
			payloadValidationErrors = append(payloadValidationErrors, payloadValidationError)
		}
		validationError := envelope.ValidationError{
			ErrorType:       &InvalidPayload.Type,
			ErrorResolution: &InvalidPayload.Resolution,
			Errors:          payloadValidationErrors,
		}
		log.Debug().Msg("event validated in " + time.Now().UTC().Sub(startTime).String())
		return false, validationError
	}
}
