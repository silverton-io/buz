package cache

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
)

const (
	GCS string = "gcs"
	S3  string = "s3"
	FS  string = "fs"
)

type SchemaCacheBackend interface {
	Initialize(config config.SchemaCacheBackend)
	GetRemote(schema string) (contents []byte, err error)
	Close()
}

func BuildSchemaCacheBackend(config config.SchemaCacheBackend) (backend SchemaCacheBackend, err error) {
	switch config.Type {
	case GCS:
		cacheBackend := GcsSchemaCacheBackend{}
		cacheBackend.Initialize(config)
		return &cacheBackend, nil
	case S3:
		cacheBackend := S3SchemaCacheBackend{}
		cacheBackend.Initialize(config)
		return &cacheBackend, nil
	case FS:
		cacheBackend := FilesystemCacheBackend{}
		cacheBackend.Initialize(config)
		return &cacheBackend, nil
	default:
		e := errors.New("unsupported schema cache backend: " + config.Type)
		log.Fatal().Err(e).Msg("unsupported backend")
		return nil, e
	}
}
