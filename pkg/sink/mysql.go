package sink

import (
	"context"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func generateMysqlDsn(conf config.Sink) string {
	port := strconv.FormatUint(uint64(conf.MysqlPort), 10)
	return conf.MysqlUser + ":" + conf.MysqlPass + "@tcp(" + conf.MysqlHost + ":" + port + ")/" + conf.MysqlDbName
}

type MysqlSink struct {
	id           *uuid.UUID
	name         string
	gormDb       *gorm.DB
	validTable   string
	invalidTable string
}

func (s *MysqlSink) Id() *uuid.UUID {
	return s.id
}

func (s *MysqlSink) Name() string {
	return s.name
}

func (s *MysqlSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing mysql sink")
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	connString := generateMysqlDsn(conf)
	gormDb, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("could not open mysql connection")
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
		} else {
			log.Debug().Msg(tbl + " table already exists - not ensuring")
		}
	}
	return nil
}

func (s *MysqlSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.gormDb.Table(s.validTable).Create(envelopes)
}

func (s *MysqlSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.gormDb.Table(s.invalidTable).Create(envelopes)
}

func (s *MysqlSink) Close() {
	log.Debug().Msg("closing mysql sink")
	db, _ := s.gormDb.DB()
	db.Close()
}
