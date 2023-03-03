// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package annotator

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/registry"
	"github.com/silverton-io/buz/pkg/validator"
	"github.com/tidwall/gjson"
)

func getMetadataFromSchema(schema []byte) envelope.EventMeta {
	schemaContents := gjson.ParseBytes(schema)
	vendor := schemaContents.Get("self.vendor").String()
	namespace := schemaContents.Get("self.namespace").String()
	version := schemaContents.Get("self.version").String()
	format := schemaContents.Get("self.format").String()
	disableValidation := schemaContents.Get("disableValidation").Bool()
	return envelope.EventMeta{
		Vendor:            vendor,
		Namespace:         namespace,
		Version:           version,
		Format:            format,
		DisableValidation: disableValidation,
	}
}

func Annotate(envelopes []envelope.Envelope, registry *registry.Registry) []envelope.Envelope {
	var e []envelope.Envelope
	for _, envelope := range envelopes {
		log.Debug().Msg("🟡 annotating event")
		// NOTE - this has the potential to be confusing in the case that
		// schema-level validation is disabled.
		// Payload validation is still executed in that case but the outcome is disregarded.
		schemaExists, validationError, schemaContents := validator.GetSchemaContents(envelope.EventMeta.Schema, registry)
		if !schemaExists {
			// If schema doesn't exist
			// consider the payload invalid as could have disabled validation if wanted
			valid := false
			envelope.Validation.IsValid = &valid
			envelope.Validation.Error = &validationError
		} else {
			m := getMetadataFromSchema(schemaContents)
			if m.Namespace != "" {
				envelope.EventMeta.Namespace = m.Namespace
			}
			envelope.EventMeta.Vendor = m.Vendor
			envelope.EventMeta.Version = m.Version
			envelope.EventMeta.Format = m.Format
			if m.DisableValidation {
				// If schema-level validation is disabled
				// consider the payload valid.
				valid := true
				envelope.Validation.IsValid = &valid
				envelope.EventMeta.DisableValidation = m.DisableValidation
			} else {
				isValid, validationError := validator.ValidatePayload(schemaContents, envelope.Payload)
				envelope.Validation.IsValid = &isValid
				if !isValid {
					// Annotate the envelope with associated validation errors
					envelope.Validation.Error = &validationError
				}
			}
		}

		e = append(e, envelope)
	}
	return e
}
