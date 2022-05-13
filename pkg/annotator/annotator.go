package annotator

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/tidwall/gjson"
)

func getMetadataFromSchema(schema []byte) envelope.EventMetadata {
	schemaContents := gjson.ParseBytes(schema)
	vendor := schemaContents.Get("self.vendor").String()
	primaryNamespace := schemaContents.Get("self.primaryNamespace").String()
	secondaryNamespace := schemaContents.Get("self.secondaryNamespace").String()
	tertiaryNamespace := schemaContents.Get("self.tertiaryNamespace").String()
	name := schemaContents.Get("self.name").String()
	version := schemaContents.Get("self.version").String()
	format := schemaContents.Get("self.format").String()
	path := schemaContents.Get("title").String()
	return envelope.EventMetadata{
		Vendor:             vendor,
		PrimaryNamespace:   primaryNamespace,
		SecondaryNamespace: secondaryNamespace,
		TertiaryNamespace:  tertiaryNamespace,
		Name:               name,
		Version:            version,
		Format:             format,
		Path:               path,
	}
}

func Annotate(envelopes []envelope.Envelope, cache *cache.SchemaCache) []envelope.Envelope {
	var e []envelope.Envelope
	for _, envelope := range envelopes {
		log.Debug().Msg("annotating event")
		switch envelope.EventMetadata.Protocol {
		case protocol.PIXEL:
			e = append(e, envelope)
		case protocol.WEBHOOK:
			e = append(e, envelope)
		case protocol.RELAY:
			e = append(e, envelope)
		default:
			// isValid, validationError, schemaContents := validator.ValidateEvent(envelope.Payload, cache)
			// envelope.ValidationMetadata.IsValid = isValid
			// envelope.ValidationMetadata.ValidationError = &validationError
			// m := getMetadataFromSchema(schemaContents)
			// envelope.EventMetadata.Vendor = m.Vendor
			// envelope.EventMetadata.PrimaryNamespace = m.PrimaryNamespace
			// envelope.EventMetadata.SecondaryNamespace = m.SecondaryNamespace
			// envelope.EventMetadata.TertiaryNamespace = m.TertiaryNamespace
			// envelope.EventMetadata.Name = m.Name
			// envelope.EventMetadata.Version = m.Version
			// envelope.EventMetadata.Path = m.Path
			e = append(e, envelope)
		}
	}
	return e
}
