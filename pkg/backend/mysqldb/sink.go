// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package mysqldb

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
	"github.com/silverton-io/buz/pkg/envelope"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	gormDb   *gorm.DB
	input    chan []envelope.Envelope
	shutdown chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing mysql sink")
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	connParams := db.ConnectionParams{
		Host: conf.Hosts[0], // Only use the first configured host
		Port: conf.Port,
		Db:   conf.Database,
		User: conf.User,
		Pass: conf.Password,
	}
	connString := generateDsn(connParams)
	gormDb, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open mysql connection")
		return err
	}
	s.gormDb = gormDb
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
	for _, tbl := range []string{s.metadata.DefaultOutput, s.metadata.DeadletterOutput} {
		ensureErr := db.EnsureTable(s.gormDb, tbl, &envelope.Envelope{})
		if ensureErr != nil {
			return ensureErr
		}
	}
	return nil
}

func (s *Sink) StartWorker() error {
	err := backendutils.StartSinkWorker(s.input, s.shutdown, s)
	return err
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.input <- envelopes
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope, output string) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	err := s.gormDb.Table(output).Create(envelopes).Error
	return err
}

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	db, _ := s.gormDb.DB()
	s.shutdown <- 1
	err := db.Close()
	return err
}
