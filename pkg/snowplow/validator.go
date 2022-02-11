package snowplow

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
)

func validateSelfDescribingPayload(payload SelfDescribingPayload, cache *cache.SchemaCache) (isValid bool, validationErrs []gojsonschema.ResultError, metadata SelfDescribingMetadata) {
	log.Debug().Msg("validating self describing payload")
	startTime := time.Now()
	schemaExists, schemaContents := cache.Get(payload.Schema)
	// Short-circuit if schema can't be found in either cache or remote backend.
	if !schemaExists {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return false, nil, SelfDescribingMetadata{}
	}
	docLoader := gojsonschema.NewGoLoader(payload.Data)
	schemaLoader := gojsonschema.NewBytesLoader(schemaContents)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not validate payload " + payload.Schema) // TODO! Test an actual invalid schema here
		return false, nil, SelfDescribingMetadata{}
	}
	contents := gjson.ParseBytes(schemaContents)
	eventMetadata := SelfDescribingMetadata{
		Event_vendor:  contents.Get("self.vendor").String(),
		Event_name:    contents.Get("self.name").String(),
		Event_format:  contents.Get("self.format").String(),
		Event_version: contents.Get("self.version").String(),
	}
	if result.Valid() {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return true, nil, eventMetadata
	} else {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return false, result.Errors(), eventMetadata
	}
}

func ValidateEvent(event Event, cache *cache.SchemaCache) (isValid bool, validationErrs []gojsonschema.ResultError, metadata SelfDescribingMetadata) {
	// Only validate if event type is self describing
	if event.Event == SELF_DESCRIBING_EVENT {
		return validateSelfDescribingPayload(SelfDescribingPayload(*event.Self_describing_event), cache)
	} else {
		return true, nil, SelfDescribingMetadata{}
	}
}

// func ValidateSelfDescribingContext() (isValid bool, err error) {

// }
