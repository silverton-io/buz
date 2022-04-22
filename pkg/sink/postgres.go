package sink

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func generatePgDsn(conf config.Sink) string {
	// postgresql://[user[:password]@][netloc][:port][/dbname]
	port := strconv.FormatUint(uint64(conf.PgPort), 10)
	return "postgresql://" + conf.PgUser + ":" + conf.PgPass + "@" + conf.PgHost + ":" + port + "/" + conf.PgDbName
}

type PostgresSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	gormDb           *gorm.DB
	validTable       string
	invalidTable     string
}

func (s *PostgresSink) Id() *uuid.UUID {
	return s.id
}

func (s *PostgresSink) Name() string {
	return s.name
}

func (s *PostgresSink) Type() string {
	return POSTGRES
}

func (s *PostgresSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *PostgresSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing postgres sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	connString := generatePgDsn(conf)
	gormDb, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("could not open pg connection")
		return err
	}
	s.gormDb = gormDb
	s.validTable, s.invalidTable = conf.ValidTable, conf.InvalidTable
	for _, tbl := range []string{s.validTable, s.invalidTable} {
		tblExists := s.gormDb.Migrator().HasTable(tbl)
		if !tblExists {
			log.Debug().Msg(tbl + " table doesn't exist - ensuring")
			err = s.gormDb.Table(tbl).AutoMigrate(&envelope.PgEnvelope{})
			if err != nil {
				log.Error().Err(err).Msg("could not auto migrate table")
				return err
			}
		} else {
			log.Debug().Msg(tbl + " table already exists - not ensuring")
		}
	}
	return nil
}

func (s *PostgresSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.validTable).Create(envelopes).Error
	return err
}

func (s *PostgresSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.invalidTable).Create(envelopes).Error
	return err
}

func (s *PostgresSink) Close() {
	log.Debug().Msg("closing postgres sink")
	db, _ := s.gormDb.DB()
	db.Close()
}
