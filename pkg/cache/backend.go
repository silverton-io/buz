package cache

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
)

const (
	GCS   string = "gcs"
	S3    string = "s3"
	FS    string = "fs"
	HTTP  string = "http"
	HTTPS string = "https"
)

type SchemaCacheBackend interface {
	Initialize(config config.SchemaCacheBackend)
	GetRemote(schema string) (contents []byte, err error)
	Close()
}

func BuildSchemaCacheBackend(conf config.SchemaCacheBackend) (backend SchemaCacheBackend, err error) {
	switch conf.Type {
	case GCS:
		cacheBackend := GcsSchemaCacheBackend{}
		cacheBackend.Initialize(conf)
		return &cacheBackend, nil
	case S3:
		cacheBackend := S3SchemaCacheBackend{}
		cacheBackend.Initialize(conf)
		return &cacheBackend, nil
	case FS:
		cacheBackend := FilesystemCacheBackend{}
		cacheBackend.Initialize(conf)
		return &cacheBackend, nil
	case HTTP:
		cacheBackend := HttpSchemaCacheBackend{}
		cacheBackend.Initialize(conf)
		return &cacheBackend, nil
	case HTTPS:
		cacheBackend := HttpSchemaCacheBackend{}
		cacheBackend.Initialize(conf)
		return &cacheBackend, nil
	default:
		e := errors.New("unsupported schema cache backend: " + conf.Type)
		log.Fatal().Err(e).Msg("unsupported backend")
		return nil, e
	}
}
