package snowplow

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/xeipuuv/gojsonschema"
)

func ValidateSelfDescribingPayload(payload SelfDescribingPayload, schemaCache *cache.SchemaCache) (bool, []gojsonschema.ResultError) {
	startTime := time.Now()
	schemaContents := schemaCache.Get(payload.Schema)
	docLoader := gojsonschema.NewGoLoader(payload.Data)
	schemaLoader := gojsonschema.NewBytesLoader(schemaContents)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not validate payload " + payload.Schema)
	}
	if result.Valid() {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return true, nil
	} else {
		log.Debug().Msg("event validated in " + time.Now().Sub(startTime).String())
		return false, result.Errors()
	}
}

// func ValidateSelfDescribingContext() (isValid bool, err error) {

// }
