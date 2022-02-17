package cache

import (
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	h "github.com/silverton-io/gosnowplow/pkg/http"
	"github.com/silverton-io/gosnowplow/pkg/util"
)

type HttpSchemaCacheBackend struct {
	host string
	path string
}

func (b *HttpSchemaCacheBackend) Initialize(conf config.SchemaCacheBackend) {
	log.Debug().Msg("initializing http schema cache backend")
	b.host = conf.Host
	b.path = conf.Path
	// Auth? TBD
}

func (b *HttpSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	schemaLocation, err := url.Parse(b.host)
	schemaLocation := url.Join(b.host, b.path, schema)
	util.PrettyPrint(schemaLocation)
	content, err := h.Get(schemaLocation)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not get schema from http schema cache backend")
		return nil, err
	}
	return content, nil
}

func (b *HttpSchemaCacheBackend) Close() {
	log.Debug().Msg("closing http schema cache backend")
	// Knock off auth tokens? TBD
}
