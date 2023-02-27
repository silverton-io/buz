// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package postgresdb

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/db"
	"github.com/silverton-io/buz/pkg/envelope"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Sink struct {
	id               *uuid.UUID
	sinkType         string
	name             string
	deliveryRequired bool
	gormDb           *gorm.DB
	validTable       string
	invalidTable     string
	inputChan        chan []envelope.Envelope
	shutdownChan     chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return backendutils.SinkMetadata{
		Id:               s.id,
		Name:             s.name,
		SinkType:         s.sinkType,
		DeliveryRequired: s.deliveryRequired,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing postgres database sink")
	id := uuid.New()
	s.id, s.sinkType, s.name, s.deliveryRequired = &id, conf.Type, conf.Name, conf.DeliveryRequired
	connParams := db.ConnectionParams{
		Host: conf.PgHost,
		Port: conf.PgPort,
		Db:   conf.PgDbName,
		User: conf.PgUser,
		Pass: conf.PgPass,
	}
	connString := GenerateDsn(connParams)
	gormDb, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open pg connection")
		return err
	}
	s.gormDb, s.validTable, s.invalidTable = gormDb, constants.BUZ_VALID_EVENTS, constants.BUZ_INVALID_EVENTS
	s.inputChan = make(chan []envelope.Envelope, 10000)
	s.shutdownChan = make(chan int, 1)
	for _, tbl := range []string{s.validTable, s.invalidTable} {
		ensureErr := db.EnsureTable(s.gormDb, tbl, &envelope.JsonbEnvelope{})
		if ensureErr != nil {
			return ensureErr
		}
	}
	return nil
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.validTable).Create(envelopes).Error // FIXME -> shard
	return err
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.inputChan <- envelopes
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	return nil
}

func (s *Sink) Shutdown() error {
	log.Debug().Msg("🟡 closing postgres database sink")
	db, _ := s.gormDb.DB()
	s.shutdownChan <- 1
	err := db.Close()
	return err
}
