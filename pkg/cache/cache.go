package cache

import (
	"github.com/coocood/freecache"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
)

type SchemaCache struct {
	Cache        *freecache.Cache
	Backend      SchemaCacheBackend
	maxSizeBytes int
	ttlSeconds   int
}

func (s *SchemaCache) Initialize(conf config.SchemaCache) error {
	cacheBackend, _ := BuildSchemaCacheBackend(conf.Backend) // FIXME - pass err up
	initErr := InitializeSchemaCacheBackend(conf.Backend, cacheBackend)
	if initErr != nil {
		return initErr
	}
	s.Backend = cacheBackend
	s.Cache = freecache.NewCache(conf.MaxSizeBytes)
	s.maxSizeBytes = conf.MaxSizeBytes
	s.ttlSeconds = conf.TtlSeconds
	log.Info().Msg("schema cache with " + conf.Type + " backend initialized")
	return nil
}

func (s *SchemaCache) Get(key string) (exists bool, data []byte) {
	k := []byte(key)
	schemaContents, _ := s.Cache.Get(k)
	if schemaContents != nil { // Schema already cached locally
		log.Debug().Msg("found cache key " + key)
		return true, schemaContents
	} else { // Schema not yet cached locally - getting from remote backend
		schemaContents, err := s.Backend.GetRemote(key)
		if err != nil { // Error when getting schema from remote backend
			log.Debug().Msg("error when getting remote schema")
			return false, nil
		}
		log.Debug().Msg("caching " + key)
		err = s.Cache.Set(k, schemaContents, s.ttlSeconds)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error when setting key " + key)
		}
		log.Debug().Msg(key + " cached successfully")
		return true, schemaContents // Schema was aquired from remote backed and cached successfully
	}
}
