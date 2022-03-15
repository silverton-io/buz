package validator

import (
	"encoding/json"

	ce "github.com/cloudevents/sdk-go/v2/event"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/silverton-io/honeypot/pkg/util"
)

type Validator struct {
	cache *cache.SchemaCache
}

func (v *Validator) ValidateEnvelope(envelope *event.Envelope) {
	switch envelope.EventProtocol {
	case protocol.SNOWPLOW:
		util.Pprint("made it here!")
		var spEvent snowplow.Event
		bytes, err := json.Marshal(envelope.Event)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not marshal envelope event")
		}
		err = json.Unmarshal(bytes, &spEvent)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not unmarshal to snowplow event")
		}
		isValid, validationError, _ := validateSnowplowEvent(spEvent, v.cache)
		envelope.IsValid = &isValid
		envelope.ValidationError = &validationError

	case protocol.CLOUDEVENTS:
		util.Pprint("made it here too!")
		var cEvent ce.Event
		bytes, err := json.Marshal(envelope.Event)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not marshal envelope event")
		}
		err = json.Unmarshal(bytes, &cEvent)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not unmarshal to cloudevent")
		}
		isValid, validationError, _ := validateCloudEvent(cEvent, v.cache)
		envelope.IsValid = &isValid
		envelope.ValidationError = &validationError
	// case protocol.GENERIC:
	default:
		isValid := false
		envelope.IsValid = &isValid
	}
}
