package validator

// import (
// 	"encoding/json"

// 	ce "github.com/cloudevents/sdk-go/v2/event"
// 	"github.com/rs/zerolog/log"
// 	"github.com/silverton-io/honeypot/pkg/cache"
// 	"github.com/silverton-io/honeypot/pkg/event"
// 	"github.com/silverton-io/honeypot/pkg/protocol"
// 	"github.com/silverton-io/honeypot/pkg/util"
// )

// type Validator struct {
// 	cache *cache.SchemaCache
// }

// func (v *Validator) ValidateEnvelope(e *event.Event) {
// 	switch e.Protocol() {
// 	case protocol.SNOWPLOW:
// 		isValid, validationError, _ := validateSnowplowEvent(e, v.cache)
// 		envelope.IsValid = &isValid
// 		envelope.ValidationError = &validationError

// 	case protocol.CLOUDEVENTS:
// 		util.Pprint("made it here too!")
// 		var cEvent ce.Event
// 		bytes, err := json.Marshal(envelope.Event)
// 		if err != nil {
// 			log.Error().Stack().Err(err).Msg("could not marshal envelope event")
// 		}
// 		err = json.Unmarshal(bytes, &cEvent)
// 		if err != nil {
// 			log.Error().Stack().Err(err).Msg("could not unmarshal to cloudevent")
// 		}
// 		isValid, validationError, _ := validateCloudEvent(cEvent, v.cache)
// 		envelope.IsValid = &isValid
// 		envelope.ValidationError = &validationError
// 	// case protocol.GENERIC:
// 	default:
// 		isValid := false
// 		envelope.IsValid = &isValid
// 	}
// }
