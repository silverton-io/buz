// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/db"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickhouseSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	gormDb           *gorm.DB
	validTable       string
	invalidTable     string
}

func (s *ClickhouseSink) Id() *uuid.UUID {
	return s.id
}

func (s *ClickhouseSink) Name() string {
	return s.name
}

func (s *ClickhouseSink) Type() string {
	return db.CLICKHOUSE
}

func (s *ClickhouseSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *ClickhouseSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing clickhouse sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	connParams := db.ConnectionParams{
		Host: conf.ClickhouseHost,
		Port: conf.ClickhousePort,
		Db:   conf.ClickhouseDbName,
		User: conf.ClickhouseUser,
		Pass: conf.ClickhousePass,
	}
	connString := db.GenerateClickhouseDsn(connParams)
	gormDb, err := gorm.Open(clickhouse.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open clickhouse connection")
		return err
	}
	s.gormDb, s.validTable, s.invalidTable = gormDb, conf.ValidTable, conf.InvalidTable
	for _, tbl := range []string{s.validTable, s.invalidTable} {
		ensureErr := db.EnsureTable(s.gormDb, tbl, &envelope.StringEnvelope{})
		if ensureErr != nil {
			return ensureErr
		}
	}
	return nil
}

func (s *ClickhouseSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.validTable).Create(envelopes).Error
	return err
}

func (s *ClickhouseSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.invalidTable).Create(envelopes).Error
	return err
}

func (s *ClickhouseSink) Close() {
	log.Debug().Msg("closing mysql sink")
	db, _ := s.gormDb.DB()
	db.Close()
}
