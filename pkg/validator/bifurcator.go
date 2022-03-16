package validator

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/event"
)

func BifurcateAndAnnotate(envelopedEvents []event.Envelope, cache *cache.SchemaCache) (validEvents []event.Envelope, invalidEvents []event.Envelope) {
	var vEvents []event.Envelope
	var invEvents []event.Envelope
	for _, envelopedEvent := range envelopedEvents {
		isValid, validationError, _ := ValidateEvent(envelopedEvent.Payload, cache)
		envelopedEvent.IsValid = &isValid
		envelopedEvent.ValidationError = &validationError
		if isValid {
			vEvents = append(vEvents, envelopedEvent)
		} else {
			invEvents = append(invEvents, envelopedEvent)
		}
	}
	return vEvents, invEvents
}
