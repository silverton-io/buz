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
	primaryNamespace := schemaContents.Get("self.primaryNamespace").String()
	secondaryNamespace := schemaContents.Get("self.secondaryNamespace").String()
	tertiaryNamespace := schemaContents.Get("self.tertiaryNamespace").String()
	name := schemaContents.Get("self.name").String()
	version := schemaContents.Get("self.version").String()
	path := schemaContents.Get("title").String()
	return envelope.EventMeta{
		Vendor:             vendor,
		PrimaryNamespace:   primaryNamespace,
		SecondaryNamespace: secondaryNamespace,
		TertiaryNamespace:  tertiaryNamespace,
		Name:               name,
		Version:            version,
		Path:               path,
	}
}

func Annotate(envelopes []envelope.Envelope, cache *cache.SchemaCache) []envelope.Envelope {
	var e []envelope.Envelope
	for _, envelope := range envelopes {
		log.Debug().Msg("annotating event")
		switch envelope.EventMeta.Protocol {
		case protocol.PIXEL:
			e = append(e, envelope)
		case protocol.WEBHOOK:
			e = append(e, envelope)
		case protocol.RELAY:
			e = append(e, envelope)
		default:
			isValid, validationError, schemaContents := validator.ValidateEvent(envelope.Payload, cache)
			envelope.Validation.IsValid = isValid
			envelope.Validation.Error = &validationError
			m := getMetadataFromSchema(schemaContents)
			envelope.EventMeta.Vendor = m.Vendor
			envelope.EventMeta.PrimaryNamespace = m.PrimaryNamespace
			envelope.EventMeta.SecondaryNamespace = m.SecondaryNamespace
			envelope.EventMeta.TertiaryNamespace = m.TertiaryNamespace
			envelope.EventMeta.Name = m.Name
			envelope.EventMeta.Version = m.Version
			envelope.EventMeta.Path = m.Path
			e = append(e, envelope)
		}
	}
	return e
}
