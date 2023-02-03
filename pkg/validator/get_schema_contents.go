// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package validator

import (
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/registry"
)

func GetSchemaContents(schemaName string, registry *registry.Registry) (exists bool, validationError envelope.ValidationError, contents []byte) {
	// FIXME- Short-circuit if the event is an unknown event
	if schemaName == "" {
		validationError := envelope.ValidationError{
			ErrorType:       &NoSchemaAssociated.Type,
			ErrorResolution: &NoSchemaAssociated.Resolution,
			Errors:          nil,
		}
		return false, validationError, nil
	}
	schemaExists, schemaContents := registry.Get(schemaName)
	if !schemaExists {
		validationError := envelope.ValidationError{
			ErrorType:       &NoSchemaInBackend.Type,
			ErrorResolution: &NoSchemaInBackend.Resolution,
			Errors:          nil,
		}
		return false, validationError, nil
	}
	return schemaExists, envelope.ValidationError{}, schemaContents
}
