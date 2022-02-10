package cache

import (
	"github.com/coocood/freecache"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
)

type SchemaCache struct {
	cache        *freecache.Cache
	backend      *GcsSchemaCacheBackend
	maxSizeBytes int
	ttlSeconds   int
}

func (s *SchemaCache) Initialize(config config.Cache) {
	switch config.Backend.Type {
	case "gcs":
		cacheBackend := GcsSchemaCacheBackend{}
		cacheBackend.Initialize(config.Backend.Location, config.Backend.Path)
		s.backend = &cacheBackend
	}
	s.cache = freecache.NewCache(config.MaxSizeBytes)
	s.maxSizeBytes = config.MaxSizeBytes
	s.ttlSeconds = config.TtlSeconds
	log.Info().Msg("schema cache initialized")
}

func (s *SchemaCache) Get(key string) (data []byte, exists bool) {
	k := []byte(key)
	schemaContents, _ := s.cache.Get(k)
	if schemaContents != nil {
		log.Debug().Msg("found cache key " + key)
		// Found schema cached locally
		return schemaContents, true
		// Did not find schema locally, going to remote
	} else {
		schemaContents, err := s.backend.getRemoteSchema(key)
		if err != nil {
			log.Debug().Msg("error when getting remote schema")
			// Can not get remote schema
			return nil, false
		}
		log.Debug().Msg("caching " + key)
		err = s.cache.Set(k, schemaContents, s.ttlSeconds)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error when setting key " + key)
		}
		log.Debug().Msg(key + " cached successfully")
		// Remote schema exists and was cached successfully
		return schemaContents, true
	}
}
