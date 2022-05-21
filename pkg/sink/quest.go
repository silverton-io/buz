package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/db"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type QuestSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	gormDb           *gorm.DB
	validTable       string
	invalidTable     string
}

func (s *QuestSink) Id() *uuid.UUID {
	return s.id
}

func (s *QuestSink) Name() string {
	return s.name
}

func (s *QuestSink) Type() string {
	return db.QUEST
}

func (s *QuestSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *QuestSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing quest sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	connParams := db.ConnectionParams{
		Host: conf.QuestHost,
		Port: conf.QuestPort,
		Db:   conf.QuestDbName,
		User: conf.QuestUser,
		Pass: conf.QuestPass,
	}
	connString := db.GeneratePostgresDsn(connParams)
	gormDb, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("could not open quest connection")
		return err
	}
	s.gormDb, s.validTable, s.invalidTable = gormDb, conf.ValidTable, conf.InvalidTable
	for _, tbl := range []string{s.validTable, s.invalidTable} {
		ensureErr := db.EnsureTable(s.gormDb, tbl, &envelope.Envelope{})
		if ensureErr != nil {
			return ensureErr
		}
	}
	return nil
}

func (s *QuestSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.validTable).Create(envelopes).Error
	return err
}

func (s *QuestSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.invalidTable).Create(envelopes).Error
	return err
}

func (s *QuestSink) Close() {
	log.Debug().Msg("closing quest sink")
	db, _ := s.gormDb.DB()
	db.Close()
}
