// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package http

import (
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/request"
)

type RegistryBackend struct {
	protocol string
	host     string
	path     string
}

func (b *RegistryBackend) Initialize(conf config.Backend) error {
	log.Debug().Msg("🟡 initializing http schema cache backend")
	b.protocol = conf.Type
	b.host = conf.Host // FIXME! String trailing / if it's present (or validate it upstream)
	b.path = conf.Path // FIXME! Strip leading / if it's present (or validate it upstream)
	return nil
}

func (b *RegistryBackend) GetRemote(schema string) (contents []byte, err error) {
	schemaLocation, _ := url.Parse(b.protocol + "://" + b.host + "/" + b.path + "/" + schema) // FIXME!! There's gotta be a better way here.
	content, err := request.Get(*schemaLocation)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not get schema from http schema cache backend")
		return nil, err
	}
	return content, nil
}

func (b *RegistryBackend) Close() {
	log.Debug().Msg("🟡 closing http schema cache backend")
	// Knock off auth tokens? TBD
}
