package cache

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/db"
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
	Initialize(config config.Backend) error
	GetRemote(schema string) (contents []byte, err error)
	Close()
}

func BuildSchemaCacheBackend(conf config.Backend) (backend SchemaCacheBackend, err error) {
	switch conf.Type {
	case GCS:
		cacheBackend := GcsSchemaCacheBackend{}
		return &cacheBackend, nil
	case S3:
		cacheBackend := S3SchemaCacheBackend{}
		return &cacheBackend, nil
	case FS:
		cacheBackend := FilesystemCacheBackend{}
		return &cacheBackend, nil
	case HTTP:
		cacheBackend := HttpSchemaCacheBackend{}
		return &cacheBackend, nil
	case HTTPS:
		cacheBackend := HttpSchemaCacheBackend{}
		return &cacheBackend, nil
	case db.POSTGRES:
		cacheBackend := PostgresSchemaCacheBackend{}
		return &cacheBackend, nil
	case db.MYSQL:
		cacheBackend := MysqlSchemaCacheBackend{}
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

func InitializeSchemaCacheBackend(conf config.Backend, b SchemaCacheBackend) error {
	err := b.Initialize(conf)
	if err != nil {
		log.Error().Err(err).Msg("could not initialize schema cache backend")
		return err
	}
	log.Info().Msg(conf.Type + " schema cache backend initialized")
	return nil
}
