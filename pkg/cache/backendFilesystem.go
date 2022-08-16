// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package cache

import (
	"io/ioutil"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
)

type FilesystemCacheBackend struct {
	path string
}

func (b *FilesystemCacheBackend) Initialize(conf config.Backend) error {
	log.Debug().Msg("🟡 initializing filesystem schema cache backend")
	b.path = conf.Path
	// No-op
	return nil
}

func (b *FilesystemCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	schemaLocation := filepath.Join(b.path, schema)
	content, err := ioutil.ReadFile(schemaLocation)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not get schema from filesystem schema cache backend: " + schemaLocation)
		return nil, err
	}
	return content, nil
}

func (b *FilesystemCacheBackend) Close() {
	log.Debug().Msg("🟡 closing filesystem schema cache backend")
	// No-op
}
