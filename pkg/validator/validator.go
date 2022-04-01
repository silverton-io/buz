package validator

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/snowplow"
)

func ValidateEvent(e event.Event, cache *cache.SchemaCache) (isValid bool, validationError envelope.ValidationError, schema []byte) {
	schemaName := e.Schema()
	eventProtocol := e.Protocol()
	// Short-circuit if the event is an unknown snowplow event
	if eventProtocol == protocol.SNOWPLOW {
		if e.(snowplow.SnowplowEvent).Event.Event == snowplow.UNKNOWN_EVENT {
			validationError := envelope.ValidationError{
				ErrorType:       &UnknownSnowplowEventType.Type,
				ErrorResolution: &UnknownSnowplowEventType.Resolution,
				Errors:          nil,
			}
			return false, validationError, nil
		}
	}
	if *schemaName == "" {
		validationError := envelope.ValidationError{
			ErrorType:       &NoSchemaAssociated.Type,
			ErrorResolution: &NoSchemaAssociated.Resolution,
			Errors:          nil,
		}
		return false, validationError, nil
	}
	schemaExists, schemaContents := cache.Get(*schemaName)
	if !schemaExists {
		validationError := envelope.ValidationError{
			ErrorType:       &NoSchemaInBackend.Type,
			ErrorResolution: &NoSchemaInBackend.Resolution,
			Errors:          nil,
		}
		return false, validationError, nil
	} else {
		payload, err := e.PayloadAsByte()
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not marshal payload")
			validationError := envelope.ValidationError{
				ErrorType:       &InvalidPayload.Type,
				ErrorResolution: &InvalidPayload.Resolution,
				Errors:          nil,
			}
			return false, validationError, nil
		}
		if payload == nil {
			validationError := envelope.ValidationError{
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
