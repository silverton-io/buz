package cache

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
)

type GcsSchemaCacheBackend struct {
	location string
	path     string
	client   *storage.Client
}

func (b *GcsSchemaCacheBackend) Initialize(location string, path string) {
	ctx := context.Background()
	log.Debug().Msg("initializing gcs schema cache backend")
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not initialize gcs schema cache backend")
	}
	b.client = client
	b.location = location
	b.path = path
}

func (b *GcsSchemaCacheBackend) getRemoteSchema(schema string) []byte {
	ctx := context.Background()
	var schemaLocation string
	if b.path == "/" {
		schemaLocation = schema
	} else {
		schemaLocation = filepath.Join(b.path, schema)
	}
	log.Debug().Msg("getting file from cache backend " + schemaLocation)
	reader, err := b.client.Bucket(b.location).Object(schemaLocation).NewReader(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not get file from cache backend " + schemaLocation)
	}
	data, _ := ioutil.ReadAll(reader)
	return data
}
