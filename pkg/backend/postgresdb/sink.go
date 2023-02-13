// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package postgresdb

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/db"
	"github.com/silverton-io/buz/pkg/envelope"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Sink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	gormDb           *gorm.DB
	validTable       string
	invalidTable     string
}

func (s *Sink) Id() *uuid.UUID {
	return s.id
}

func (s *Sink) Name() string {
	return s.name
}

func (s *Sink) Type() string {
	return "postgresdb"
}

func (s *Sink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing postgres sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
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
	for _, tbl := range []string{s.validTable, s.invalidTable} {
		ensureErr := db.EnsureTable(s.gormDb, tbl, &envelope.JsonbEnvelope{})
		if ensureErr != nil {
			return ensureErr
		}
	}
	return nil
}

func (s *Sink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.validTable).Create(envelopes).Error
	return err
}

func (s *Sink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.invalidTable).Create(envelopes).Error
	return err
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	return nil
}

func (s *Sink) Close() {
	log.Debug().Msg("🟡 closing postgres sink")
	db, _ := s.gormDb.DB()
	db.Close()
}
