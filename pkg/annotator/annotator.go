// Copyright (c) 2023 Silverton Data, Inc.
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

type schemaMetadata struct {
	Vendor            string
	Namespace         string
	Version           string
	DisableValidation bool
}

func getSchemaMetadata(schema []byte) schemaMetadata {
	schemaContents := gjson.ParseBytes(schema)
	vendor := schemaContents.Get("self.vendor").String()
	namespace := schemaContents.Get("self.namespace").String()
	version := schemaContents.Get("self.version").String()
	disableValidation := schemaContents.Get("disableValidation").Bool()
	log.Debug().Msgf("🟡 vendor: %s, namespace: %s, version: %s, disableValidation: %t", vendor, namespace, version, disableValidation)
	return schemaMetadata{
		Vendor:            vendor,
		Namespace:         namespace,
		Version:           version,
		DisableValidation: disableValidation,
	}
}

func Annotate(envelopes []envelope.Envelope, registry *registry.Registry) []envelope.Envelope {
	var e []envelope.Envelope
	for _, envelope := range envelopes {
		log.Debug().Msg("🟡 annotating event")
		isValid, validationError, schemaContents := validator.Validate(envelope, registry)
		m := getSchemaMetadata(schemaContents)
		if m.Namespace != "" {
			envelope.Vendor = m.Vendor
			envelope.Namespace = m.Namespace
			envelope.Version = m.Version
		}
		if m.DisableValidation {
			// If schema-level validation is disabled
			// consider the payload valid.
			valid := true
			envelope.IsValid = valid
		} else {
			envelope.IsValid = isValid
			if !isValid {
				// Annotate the envelope with associated validation errors
				envelope.ValidationError = &validationError
			}
		}
		e = append(e, envelope)
	}
	return e
}
