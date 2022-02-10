package snowplow

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/xeipuuv/gojsonschema"
)

func validateSelfDescribingPayload(payload SelfDescribingPayload, cache *cache.SchemaCache) (bool, []gojsonschema.ResultError) {
	log.Debug().Msg("validating self describing payload")
	startTime := time.Now()
	schemaContents, schemaExists := cache.Get(payload.Schema)
	if !schemaExists {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return false, nil
	}
	docLoader := gojsonschema.NewGoLoader(payload.Data)
	schemaLoader := gojsonschema.NewBytesLoader(schemaContents)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not validate payload " + payload.Schema) // FIXME! Test an actual invalid schema here
	}
	if result.Valid() {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return true, nil
	} else {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return false, result.Errors()
	}
}

func ValidateEvent(event Event, cache *cache.SchemaCache) (bool, []gojsonschema.ResultError) {
	// Only validate if event type is self describing
	if event.Event == SELF_DESCRIBING_EVENT {
		return validateSelfDescribingPayload(SelfDescribingPayload(*event.Self_describing_event), cache)
	} else {
		return true, nil
	}
}

// func ValidateSelfDescribingContext() (isValid bool, err error) {

// }
