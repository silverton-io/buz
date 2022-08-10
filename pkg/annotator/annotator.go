// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package annotator

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func getMetadataFromSchema(schema []byte) envelope.EventMeta {
	schemaContents := gjson.ParseBytes(schema)
	vendor := schemaContents.Get("self.vendor").String()
	namespace := schemaContents.Get("self.namespace").String()
	version := schemaContents.Get("self.version").String()
	format := schemaContents.Get("self.format").String()
	return envelope.EventMeta{
		Vendor:    vendor,
		Namespace: namespace,
		Version:   version,
		Format:    format,
	}
}

func Annotate(envelopes []envelope.Envelope, cache *cache.SchemaCache) []envelope.Envelope {
	var e []envelope.Envelope
	for _, envelope := range envelopes {
		log.Debug().Msg("annotating event")
		isValid, validationError, schemaContents := validator.ValidatePayload(envelope.EventMeta.Schema, envelope.Payload, cache)
		envelope.Validation.IsValid = isValid
		envelope.Validation.Error = &validationError
		m := getMetadataFromSchema(schemaContents)
		if m.Namespace != "" {
			envelope.EventMeta.Namespace = m.Namespace
		}
		envelope.EventMeta.Vendor = m.Vendor
		envelope.EventMeta.Version = m.Version
		envelope.EventMeta.Format = m.Format
		e = append(e, envelope)
	}
	return e
}
