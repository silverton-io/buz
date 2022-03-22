package validator

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/tele"
)

func BifurcateAndAnnotate(envelopes []envelope.Envelope, cache *cache.SchemaCache) (validEvents []envelope.Envelope, invalidEvents []envelope.Envelope, stats tele.ProtocolStats) {
	var vEnvelopes []envelope.Envelope
	var invEnvelopes []envelope.Envelope
	protocolStats := tele.ProtocolStats{}
	protocolStats.Build()
	for _, envelope := range envelopes {
		isValid, validationError, _ := ValidateEvent(envelope.Payload, cache)
		envelope.IsValid = &isValid
		envelope.ValidationError = &validationError
		// FIXME - Annotate envelope root with schema information.
		if isValid {
			vEnvelopes = append(vEnvelopes, envelope)
			protocolStats.IncrementValid(envelope.EventProtocol, envelope.EventSchema, 1)
		} else {
			invEnvelopes = append(invEnvelopes, envelope)
			protocolStats.IncrementInvalid(envelope.EventProtocol, envelope.EventSchema, 1)
		}
	}
	return vEnvelopes, invEnvelopes, protocolStats
}

func Bifurcate(envelopes []envelope.Envelope) (validEvents []envelope.Envelope, invalidEvents []envelope.Envelope, stats tele.ProtocolStats) {
	var vEnvelopes []envelope.Envelope
	var invEnvelopes []envelope.Envelope
	protocolStats := tele.ProtocolStats{}
	protocolStats.Build()
	for _, envelope := range envelopes {
		isValid := *envelope.IsValid
		if isValid {
			vEnvelopes = append(vEnvelopes, envelope)
			protocolStats.IncrementValid(envelope.EventProtocol, envelope.EventSchema, 1)
		} else {
			invEnvelopes = append(invEnvelopes, envelope)
			protocolStats.IncrementInvalid(envelope.EventProtocol, envelope.EventSchema, 1)
		}
	}
	return vEnvelopes, invEnvelopes, protocolStats
}
