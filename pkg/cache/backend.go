package cache

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
)

const (
	GCS   string = "gcs"
	S3    string = "s3"
	FS    string = "fs"
	HTTP  string = "http"
	HTTPS string = "https"
	IGLU  string = "iglu"
	KSR   string = "ksr" // Kafka schema registry
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
	case IGLU:
		e := errors.New("the iglu schema cache backend is not yet available")
		log.Fatal().Stack().Err(e).Msg("iglu is unsupported")
		return nil, e
	case KSR:
		e := errors.New("the kafka schema registry cache backend is not yet available")
		log.Fatal().Stack().Err(e).Msg("kafka schema registry is unsupported")
		return nil, e
	default:
		e := errors.New("unsupported schema cache backend: " + conf.Type)
		log.Fatal().Stack().Err(e).Msg("unsupported backend")
		return nil, e
	}
}
