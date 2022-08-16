package cache

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
)

type MinioSchemaCacheBackend struct {
	endpoint        string
	accessKeyId     string
	secretAccessKey string
	// TODO: add support for token
	// TODO: add support for useSsl
	bucket string
	path   string
	client *minio.Client
}

func (b *MinioSchemaCacheBackend) Initialize(conf config.Backend) error {
	log.Debug().Msg("initializing minio schema cache backend")
	b.endpoint = conf.MinioEndpoint
	b.accessKeyId = conf.AccessKeyId
	b.secretAccessKey = conf.SecretAccessKey
	b.bucket, b.path = conf.Bucket, conf.Path
	client, err := minio.New(b.endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(b.accessKeyId, b.secretAccessKey, ""),
	})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not initialize minio client")
		return err
	}
	b.client = client
	return nil
}

func (b *MinioSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	ctx := context.Background()
	var schemaLocation string
	if b.path == "/" {
		schemaLocation = schema
	} else {
		schemaLocation = filepath.Join(b.path, schema)
	}
	log.Debug().Msg("getting file from minio backend " + schemaLocation)
	obj, err := b.client.GetObject(ctx, b.bucket, schemaLocation, minio.GetObjectOptions{})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not get file from minio: " + schemaLocation)
		return nil, err
	}
	contents, err = ioutil.ReadAll(obj)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not read contents from file: " + schemaLocation)
		return nil, err
	}
	return contents, nil
}

func (b *MinioSchemaCacheBackend) Close() {
	log.Debug().Msg("closing minio schema cache backend")
}
