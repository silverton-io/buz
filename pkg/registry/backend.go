// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package registry

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/clickhousedb"
	"github.com/silverton-io/buz/pkg/backend/file"
	"github.com/silverton-io/buz/pkg/backend/http"
	"github.com/silverton-io/buz/pkg/backend/materializedb"
	"github.com/silverton-io/buz/pkg/backend/mongodb"
	"github.com/silverton-io/buz/pkg/backend/mysqldb"
	"github.com/silverton-io/buz/pkg/backend/postgresdb"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
)

const (
	GCS   string = "gcs"
	S3    string = "s3"
	MINIO string = "minio"
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
	case MINIO:
		cacheBackend := MinioSchemaCacheBackend{}
		return &cacheBackend, nil
	case FS:
		cacheBackend := file.RegistryBackend{}
		return &cacheBackend, nil
	case HTTP:
		cacheBackend := http.RegistryBackend{}
		return &cacheBackend, nil
	case HTTPS:
		cacheBackend := http.RegistryBackend{}
		return &cacheBackend, nil
	case constants.POSTGRES:
		cacheBackend := postgresdb.RegistryBackend{}
		return &cacheBackend, nil
	case constants.MYSQL:
		cacheBackend := mysqldb.RegistryBackend{}
		return &cacheBackend, nil
	case constants.MATERIALIZE:
		cacheBackend := materializedb.RegistryBackend{}
		return &cacheBackend, nil
	case constants.CLICKHOUSE:
		cacheBackend := clickhousedb.RegistryBackend{}
		return &cacheBackend, nil
	case constants.MONGODB:
		cacheBackend := mongodb.RegistryBackend{}
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
		log.Fatal().Stack().Err(e).Msg("🔴 unsupported backend")
		return nil, e
	}
}

func InitializeSchemaCacheBackend(conf config.Backend, b SchemaCacheBackend) error {
	err := b.Initialize(conf)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not initialize schema cache backend")
		return err
	}
	log.Info().Msg("🟢 " + conf.Type + " schema cache backend initialized")
	return nil
}
