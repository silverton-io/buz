package validator

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/snowplow"
)

func ValidateEvent(e event.Event, cache *cache.SchemaCache) (isValid bool, validationError event.ValidationError, schema []byte) {
	schemaName := e.Schema()
	eventProtocol := e.Protocol()
	// Short-circuit if the event is an unknown snowplow event or if it is a snowplow event but not self-describing
	if eventProtocol == protocol.SNOWPLOW {
		if e.(snowplow.SnowplowEvent).Event == snowplow.UNKNOWN_EVENT {
			validationError := event.ValidationError{
				ErrorType:       &UnknownSnowplowEventType.Type,
				ErrorResolution: &UnknownSnowplowEventType.Resolution,
				Errors:          nil,
			}
			return false, validationError, nil
		}
		if e.(snowplow.SnowplowEvent).Event != snowplow.SELF_DESCRIBING_EVENT {
			return true, event.ValidationError{}, nil
		}
	}
	if *schemaName == "" {
		validationError := event.ValidationError{
			ErrorType:       &NoSchemaAssociated.Type,
			ErrorResolution: &NoSchemaAssociated.Resolution,
			Errors:          nil,
		}
		return false, validationError, nil
	}
	schemaExists, schemaContents := cache.Get(*schemaName)
	if !schemaExists {
		validationError := event.ValidationError{
			ErrorType:       &NoSchemaInBackend.Type,
			ErrorResolution: &NoSchemaInBackend.Resolution,
			Errors:          nil,
		}
		return false, validationError, nil
	} else {
		payload, err := e.PayloadAsByte()
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not marshal payload")
			validationError := event.ValidationError{
				ErrorType:       &InvalidPayload.Type,
				ErrorResolution: &InvalidPayload.Resolution,
				Errors:          nil,
			}
			return false, validationError, nil
		}
		if payload == nil {
			validationError := event.ValidationError{
				ErrorType:       &PayloadNotPresent.Type,
				ErrorResolution: &PayloadNotPresent.Resolution,
				Errors:          nil,
			}
			return false, validationError, nil
		}
		isValid, validationError := validatePayload(payload, schemaContents)
		return isValid, validationError, schemaContents
	}
}
