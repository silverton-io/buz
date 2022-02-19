package cache

import (
	"github.com/coocood/freecache"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
)

type SchemaCache struct {
	cache        *freecache.Cache
	Backend      SchemaCacheBackend
	maxSizeBytes int
	ttlSeconds   int
}

func (s *SchemaCache) Initialize(config config.SchemaCache) {
	cacheBackend, _ := BuildSchemaCacheBackend(config.SchemaCacheBackend)
	s.Backend = cacheBackend
	s.cache = freecache.NewCache(config.MaxSizeBytes)
	s.maxSizeBytes = config.MaxSizeBytes
	s.ttlSeconds = config.TtlSeconds
	log.Info().Msg("schema cache with " + config.Type + " backend initialized")
}

func (s *SchemaCache) Get(key string) (exists bool, data []byte) {
	k := []byte(key)
	schemaContents, _ := s.cache.Get(k)
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
		err = s.cache.Set(k, schemaContents, s.ttlSeconds)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error when setting key " + key)
		}
		log.Debug().Msg(key + " cached successfully")
		return true, schemaContents // Schema was aquired from remote backed and cached successfully
	}
}
