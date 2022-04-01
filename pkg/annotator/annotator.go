package annotator

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func getMetadata(schema []byte) envelope.EventMetadata {
	schemaContents := gjson.ParseBytes(schema)
	vendor := schemaContents.Get("self.vendor").String()
	primaryCategory := schemaContents.Get("self.primaryCategory").String()
	secondaryCategory := schemaContents.Get("self.secondaryCategory").String()
	tertiaryCategory := schemaContents.Get("self.tertiaryCategory").String()
	name := schemaContents.Get("self.name").String()
	version := schemaContents.Get("self.version").String()
	format := schemaContents.Get("self.format").String()
	path := schemaContents.Get("title").String()
	return envelope.EventMetadata{
		Vendor:            &vendor,
		PrimaryCategory:   &primaryCategory,
		SecondaryCategory: &secondaryCategory,
		TertiaryCategory:  &tertiaryCategory,
		Name:              &name,
		Version:           &version,
		Format:            &format,
		Path:              &path,
	}
}

func Annotate(envelopes []envelope.Envelope, cache *cache.SchemaCache) []envelope.Envelope {
	var e []envelope.Envelope
	for _, envelope := range envelopes {
		isValid, validationError, schemaContents := validator.ValidateEvent(envelope.Payload, cache)
		envelope.IsValid = &isValid
		envelope.ValidationError = &validationError
		envelope.EventMetadata = getMetadata(schemaContents)
		e = append(e, envelope)
	}
	return e
}
