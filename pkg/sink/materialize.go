// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package sink

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/db"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func generateMzDsn(conf config.Sink) string {
	// postgresql://[user[:password]@][netloc][:port][/dbname]
	port := strconv.FormatUint(uint64(conf.MzPort), 10)
	return "postgresql://" + conf.MzUser + ":" + conf.MzPass + "@" + conf.MzHost + ":" + port + "/" + conf.MzDbName
}

type MaterializeSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	gormDb           *gorm.DB
	validTable       string
	invalidTable     string
}

func (s *MaterializeSink) Id() *uuid.UUID {
	return s.id
}

func (s *MaterializeSink) Name() string {
	return s.name
}

func (s *MaterializeSink) Type() string {
	return db.MATERIALIZE
}

func (s *MaterializeSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *MaterializeSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing materialize sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	connString := generateMzDsn(conf)
	gormDb, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("could not open materialize connection")
		return err
	}
	s.gormDb = gormDb
	s.validTable, s.invalidTable = conf.ValidTable, conf.InvalidTable
	for _, tbl := range []string{s.validTable, s.invalidTable} {
		ensureErr := db.EnsureTable(s.gormDb, tbl, &envelope.Envelope{})
		if ensureErr != nil {
			return ensureErr
		}
	}
	return nil
}

func (s *MaterializeSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.validTable).Create(envelopes).Error
	return err
}

func (s *MaterializeSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.invalidTable).Create(envelopes).Error
	return err
}

func (s *MaterializeSink) Close() {
	log.Debug().Msg("closing postgres sink")
	db, _ := s.gormDb.DB()
	db.Close()
}
