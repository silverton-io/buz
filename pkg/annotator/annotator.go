package annotator

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func getMetadataFromSchema(schema []byte) envelope.EventMeta {
	schemaContents := gjson.ParseBytes(schema)
	vendor := schemaContents.Get("self.vendor").String()
	namespace := schemaContents.Get("self.namespace").String()
	version := schemaContents.Get("self.version").String()
	format := schemaContents.Get("self.format").String()
	path := schemaContents.Get("title").String()
	return envelope.EventMeta{
		Vendor:    vendor,
		Namespace: namespace,
		Version:   version,
		Format:    format,
		Path:      path,
	}
}

func Annotate(envelopes []envelope.Envelope, cache *cache.SchemaCache) []envelope.Envelope {
	var e []envelope.Envelope
	for _, envelope := range envelopes {
		log.Debug().Msg("annotating event")
		switch envelope.EventMeta.Protocol {
		case protocol.RELAY:
			e = append(e, envelope)
		default:
			isValid, validationError, schemaContents := validator.ValidateEvent(envelope.Payload, cache)
			envelope.Validation.IsValid = isValid
			envelope.Validation.Error = &validationError
			m := getMetadataFromSchema(schemaContents)
			envelope.EventMeta.Vendor = m.Vendor
			envelope.EventMeta.Namespace = m.Namespace
			envelope.EventMeta.Version = m.Version
			envelope.EventMeta.Format = m.Format
			envelope.EventMeta.Path = m.Path
			e = append(e, envelope)
		}
	}
	return e
}
