package validator

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

func BifurcateAndAnnotate(envelopes []envelope.Envelope, cache *cache.SchemaCache) (validEvents []envelope.Envelope, invalidEvents []envelope.Envelope) {
	var vEnvelopes []envelope.Envelope
	var invEnvelopes []envelope.Envelope
	for _, envelopedEvent := range envelopes {
		isValid, validationError, _ := ValidateEvent(envelopedEvent.Payload, cache)
		envelopedEvent.IsValid = &isValid
		envelopedEvent.ValidationError = &validationError
		// FIXME - Annotate envelope root with schema information.
		if isValid {
			vEnvelopes = append(vEnvelopes, envelopedEvent)
		} else {
			invEnvelopes = append(invEnvelopes, envelopedEvent)
		}
	}
	return vEnvelopes, invEnvelopes
}

func Bifurcate(envelopes []envelope.Envelope) (validEvents []envelope.Envelope, invalidEvents []envelope.Envelope) {
	var vEnvelopes []envelope.Envelope
	var invEnvelopes []envelope.Envelope
	for _, envelope := range envelopes {
		isValid := *envelope.IsValid
		if isValid {
			vEnvelopes = append(vEnvelopes, envelope)
		} else {
			invEnvelopes = append(invEnvelopes, envelope)
		}
	}
	return vEnvelopes, invEnvelopes
}
