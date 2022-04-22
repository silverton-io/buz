package annotator

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/validator"
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
		switch envelope.EventProtocol {
		case protocol.PIXEL:
			var schema []byte
			eventMetadata := getMetadataFromSchema(schema)
			envelope.EventMetadata = &eventMetadata
			e = append(e, envelope)
		case protocol.WEBHOOK:
			var schema []byte
			eventMetadata := getMetadataFromSchema(schema)
			envelope.EventMetadata = &eventMetadata
			e = append(e, envelope)
		case protocol.RELAY:
			e = append(e, envelope) // Don't annotate
		default:
			isValid, validationError, schemaContents := validator.ValidateEvent(envelope.Payload, cache)
			envelope.IsValid = &isValid
			envelope.ValidationError = &validationError
			eventMetadata := getMetadataFromSchema(schemaContents)
			envelope.EventMetadata = &eventMetadata
			e = append(e, envelope)
		}
	}
	return e
}
