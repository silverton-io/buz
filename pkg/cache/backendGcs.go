package cache

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
)

type GcsSchemaCacheBackend struct {
	bucket string
	path   string
	client *storage.Client
}

func (b *GcsSchemaCacheBackend) Initialize(config config.Backend) error {
	ctx := context.Background()
	log.Debug().Msg("initializing gcs schema cache backend")
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not initialize gcs schema cache backend")
		return err
	}
	b.client, b.bucket, b.path = client, config.Bucket, config.Path
	return nil
}

func (b *GcsSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	ctx := context.Background()
	var schemaLocation string
	if b.path == "/" {
		schemaLocation = schema
	} else {
		schemaLocation = filepath.Join(b.path, schema)
	}
	log.Debug().Msg("getting file from gcs backend " + schemaLocation)
	reader, err := b.client.Bucket(b.bucket).Object(schemaLocation).NewReader(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not get file from gcs: " + schemaLocation)
		return nil, err
	}
	data, _ := ioutil.ReadAll(reader)
	return data, nil
}

func (b *GcsSchemaCacheBackend) Close() {
	log.Debug().Msg("closing gcs schema cache backend")
	b.client.Close()
}
