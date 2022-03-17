package validator

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/event"
)

func BifurcateAndAnnotate(envelopes []event.Envelope, cache *cache.SchemaCache) (validEvents []event.Envelope, invalidEvents []event.Envelope) {
	var vEnvelopes []event.Envelope
	var invEnvelopes []event.Envelope
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

func Bifurcate(envelopes []event.Envelope) (validEvents []event.Envelope, invalidEvents []event.Envelope) {
	var vEnvelopes []event.Envelope
	var invEnvelopes []event.Envelope
	for _, envelope := range envelopes {
		if *envelope.IsValid {
			vEnvelopes = append(vEnvelopes, envelope)
		} else {
			invEnvelopes = append(invEnvelopes, envelope)
		}
	}
	return vEnvelopes, invEnvelopes
}
