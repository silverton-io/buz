// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package clickhousedb

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
	"github.com/silverton-io/buz/pkg/envelope"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	gormDb   *gorm.DB
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing clickhouse sink")
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	connParams := db.ConnectionParams{
		Host: conf.Hosts[0], // Only use the first configured host
		Port: conf.Port,
		Db:   conf.Name,
		User: conf.User,
		Pass: conf.Password,
	}
	connString := generateDsn(connParams)
	gormDb, err := gorm.Open(clickhouse.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open clickhouse connection")
		return err
	}
	s.gormDb = gormDb
	for _, tbl := range []string{s.metadata.DefaultOutput, s.metadata.DeadletterOutput} {
		ensureErr := db.EnsureTable(s.gormDb, tbl, &envelope.StringEnvelope{})
		if ensureErr != nil {
			return ensureErr
		}
	}
	return nil
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	// Get shards
	// err := s.gormDb.Table(s.validTable).Create(envelopes).Error
	// return err
	return nil
}

func (s *Sink) Close() {
	log.Debug().Msg("🟡 closing mysql sink")
	db, _ := s.gormDb.DB()
	db.Close()
}
