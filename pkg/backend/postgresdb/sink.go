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
	id                 *uuid.UUID
	sinkType           string
	name               string
	deliveryRequired   bool
	gormDb             *gorm.DB
	defaultEventsTable string
	input              chan []envelope.Envelope
	shutdown           chan int
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
	id := uuid.New()
	s.id, s.sinkType, s.name, s.deliveryRequired = &id, conf.Type, conf.Name, conf.DeliveryRequired
	log.Debug().Msg("🟢 initializing " + s.sinkType + " sink")
	connParams := db.ConnectionParams{
		Host: conf.Hosts[0], //Only use the first configured host.
		Port: conf.Port,
		Db:   conf.Database,
		User: conf.User,
		Pass: conf.Password,
	}
	connString := GenerateDsn(connParams)
	gormDb, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open " + s.sinkType + " connection")
		return err
	}
	s.gormDb, s.defaultEventsTable = gormDb, constants.BUZ_EVENTS
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
	for _, tbl := range []string{s.defaultEventsTable} {
		ensureErr := db.EnsureTable(s.gormDb, tbl, &envelope.JsonbEnvelope{})
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

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	err := s.gormDb.Table(s.defaultEventsTable).Create(envelopes).Error
	return err
}

func (s *Sink) Shutdown() error {
	log.Debug().Msg("🟢 shutting down " + s.sinkType + " sink")
	db, _ := s.gormDb.DB()
	s.shutdown <- 1
	err := db.Close()
	return err
}
