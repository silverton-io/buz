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

func (s *SchemaCache) Get(key string) []byte {
	k := []byte(key)
	val, err := s.cache.Get(k)
	if err != nil {
		log.Debug().Msg("could not find cache key " + key)
	}
	if val != nil {
		log.Debug().Msg("found cache key " + key)
		return val
	} else {
		schemaContents := s.backend.getRemoteSchema(key)
		log.Debug().Msg("caching " + key)
		err := s.cache.Set(k, schemaContents, s.ttlSeconds)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error when setting key " + key)
		}
		log.Debug().Msg(key + " cached successfully")
		return schemaContents
	}
}
