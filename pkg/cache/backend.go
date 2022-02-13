package cache

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/util"
)

const (
	GCS string = "gcs"
	S3  string = "s3"
)

type SchemaCacheBackend interface {
	Initialize(config config.SchemaCacheBackend)
	GetRemote(schema string) (contents []byte, err error)
}

func BuildSchemaCacheBackend(config config.SchemaCacheBackend) (backend SchemaCacheBackend, err error) {
	util.PrettyPrint(config)
	switch config.Type {
	case GCS:
		cacheBackend := GcsSchemaCacheBackend{}
		cacheBackend.Initialize(config)
		return &cacheBackend, nil
	case S3:
		log.Fatal().Msg("S3 cache backend is currently unsupported.")
		return nil, nil
	default:
		log.Fatal().Msg("Unsupported cache backend")
		return nil, nil
	}
}
