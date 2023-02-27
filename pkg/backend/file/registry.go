// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package file

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
)

type RegistryBackend struct {
	path string
}

func (b *RegistryBackend) Initialize(conf config.Backend) error {
	log.Debug().Msg("🟡 initializing filesystem registry backend")
	b.path = conf.Path
	// No-op
	return nil
}

func (b *RegistryBackend) GetRemote(schema string) (contents []byte, err error) {
	schemaLocation := filepath.Join(b.path, schema)
	content, err := os.ReadFile(schemaLocation)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not get schema from filesystem registry backend: " + schemaLocation)
		return nil, err
	}
	return content, nil
}

func (b *RegistryBackend) Close() {
	log.Debug().Msg("🟡 closing filesystem registry backend")
	// No-op
}
