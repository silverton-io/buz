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
		isValid, validationError, schemaContents := validator.ValidatePayload(envelope.EventMeta.Schema, envelope.Payload, registry)
		m := getMetadataFromSchema(schemaContents)
		if m.Namespace != "" {
			envelope.EventMeta.Namespace = m.Namespace
		}
		envelope.EventMeta.Vendor = m.Vendor
		envelope.EventMeta.Version = m.Version
		envelope.EventMeta.Format = m.Format
		envelope.EventMeta.DisableValidation = m.DisableValidation
		if m.DisableValidation {
			// If schema-level validation override is in place, treat
			// the payload as valid. Regardless if that is actually the case.
			valid := true
			envelope.Validation.IsValid = &valid
			e = append(e, envelope)
		} else {
			envelope.Validation.IsValid = &isValid
			if !isValid {
				envelope.Validation.Error = &validationError
			}
			e = append(e, envelope)
		}
	}
	return e
}
