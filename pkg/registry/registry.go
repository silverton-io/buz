// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package registry

import (
	"github.com/coocood/freecache"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
)

type Registry struct {
	Cache        *freecache.Cache
	Backend      SchemaCacheBackend
	maxSizeBytes int
	ttlSeconds   int
}

func (r *Registry) Initialize(conf config.Registry) error {
	cacheBackend, _ := BuildSchemaCacheBackend(conf.Backend) // FIXME - pass err up
	initErr := InitializeSchemaCacheBackend(conf.Backend, cacheBackend)
	if initErr != nil {
		return initErr
	}
	r.Backend = cacheBackend
	r.Cache = freecache.NewCache(conf.MaxSizeBytes)
	r.maxSizeBytes = conf.MaxSizeBytes
	r.ttlSeconds = conf.TtlSeconds
	return nil
}

func (r *Registry) Get(key string) (exists bool, data []byte) {
	k := []byte(key)
	schemaContents, _ := r.Cache.Get(k)
	if schemaContents != nil { // Schema already cached locally
		log.Debug().Msg("🟡 found cache key " + key)
		return true, schemaContents
	} else { // Schema not yet cached locally - getting from remote backend
		schemaContents, err := r.Backend.GetRemote(key)
		if err != nil { // Error when getting schema from remote backend
			log.Debug().Msg("error when getting remote schema")
			return false, nil
		}
		log.Debug().Msg("🟡 caching " + key)
		err = r.Cache.Set(k, schemaContents, r.ttlSeconds)
		if err != nil {
			log.Error().Err(err).Msg("🔴 error when setting key " + key)
		}
		log.Debug().Msg("🟡 " + key + " cached successfully")
		return true, schemaContents // Schema was aquired from remote backed and cached successfully
	}
}
