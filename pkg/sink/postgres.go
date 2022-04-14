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
	port := strconv.FormatUint(uint64(conf.DbPort), 10)
	return "postgresql://" + conf.DbUser + ":" + conf.DbPass + "@" + conf.DbHost + ":" + port + "/" + conf.DbName
}

type PostgresSink struct {
	id           *uuid.UUID
	name         string
	gormDb       *gorm.DB
	validTable   string
	invalidTable string
}

func (s *PostgresSink) Id() *uuid.UUID {
	return s.id
}

func (s *PostgresSink) Name() string {
	return s.name
}

func (s *PostgresSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing postgres sink")
	id := uuid.New()
	s.id, s.name = &id, conf.Name
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
			err = s.gormDb.Table(tbl).AutoMigrate(&envelope.Envelope{})
			if err != nil {
				log.Error().Err(err).Msg("could not auto migrate table")
				return err
			}
			// NOTE! This is a hacky workaround so that the same gorm struct tag of "json" can be used, but "jsonb" is used for pg.
			for _, col := range []string{"event_metadata", "validation_error", "payload"} {
				alterStmt := "alter table " + tbl + " alter column " + col + " set data type jsonb using " + col + "::jsonb;"
				log.Debug().Msg("ensuring jsonb columns via: " + alterStmt)
				s.gormDb.Exec(alterStmt)
			}
		} else {
			log.Debug().Msg(tbl + " table already exists - not ensuring")
		}
	}
	return nil
}

func (s *PostgresSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.gormDb.Table(s.validTable).Create(envelopes)
}

func (s *PostgresSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.gormDb.Table(s.invalidTable).Create(envelopes)
}

func (s *PostgresSink) Close() {
	log.Debug().Msg("closing postgres sink")
	db, _ := s.gormDb.DB()
	db.Close()
}
