package validator

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

func Annotate(envelopes []envelope.Envelope, cache *cache.SchemaCache) []envelope.Envelope {
	var e []envelope.Envelope
	for _, envelope := range envelopes {
		isValid, validationError, _ := ValidateEvent(envelope.Payload, cache)
		envelope.IsValid = &isValid
		envelope.ValidationError = &validationError
		// FIXME - Annotate envelope root with schema information.
		e = append(e, envelope)
	}
	return e
}
