// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package registry

import (
	"strings"

	"github.com/coocood/freecache"
	"github.com/rs/zerolog/log"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/silverton-io/buz/pkg/config"
)

type Registry struct {
	Cache        *freecache.Cache
	Backend      SchemaCacheBackend
	Compiler     *jsonschema.Compiler
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
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft2020
	r.Compiler = compiler
	return nil
}

func (r *Registry) Get(key string) (exists bool, data []byte) {
	k := []byte(key)
	schemaContents, _ := r.Cache.Get(k)
	if schemaContents != nil { // Schema already cached locally
		log.Debug().Msg("🟡 found cache key " + key)
		return true, schemaContents
	} else { // Schema not yet cached locally - getting from remote backend
		// Ensure schemaKey is key ending in .json (add if not present)
		schemaKey := key
		if !strings.HasSuffix(schemaKey, ".json") {
			schemaKey = schemaKey + ".json"
		}
		schemaContents, err := r.Backend.GetRemote(schemaKey)
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
		log.Debug().Msg("🟡 adding schema to compiler " + key)
		err = r.Compiler.AddResource(key, strings.NewReader(string(schemaContents)))
		if err != nil {
			log.Error().Err(err).Msg("🔴 error when compiling schema " + key)
		}
		return true, schemaContents // Schema was aquired from remote backed and cached successfully
	}
}
