package sink

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/db"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	gormDb           *gorm.DB
	validTable       string
	invalidTable     string
}

func (s *MysqlSink) Id() *uuid.UUID {
	return s.id
}

func (s *MysqlSink) Name() string {
	return s.name
}

func (s *MysqlSink) Type() string {
	return db.MYSQL
}

func (s *MysqlSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *MysqlSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing mysql sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	connParams := db.ConnectionParams{
		Host: conf.MysqlHost,
		Port: conf.MysqlPort,
		Db:   conf.MysqlDbName,
		User: conf.MysqlUser,
		Pass: conf.MysqlPass,
	}
	connString := db.GenerateMysqlDsn(connParams)
	gormDb, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("could not open mysql connection")
		return err
	}
	s.gormDb, s.validTable, s.invalidTable = gormDb, conf.ValidTable, conf.InvalidTable
	for _, tbl := range []string{s.validTable, s.invalidTable} {
		ensureErr := db.EnsureTable(s.gormDb, tbl, &envelope.MysqlEnvelope{})
		if ensureErr != nil {
			return ensureErr
		}
	}
	return nil
}

func (s *MysqlSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.validTable).Create(envelopes).Error
	return err
}

func (s *MysqlSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.gormDb.Table(s.invalidTable).Create(envelopes).Error
	return err
}

func (s *MysqlSink) Close() {
	log.Debug().Msg("closing mysql sink")
	db, _ := s.gormDb.DB()
	db.Close()
}
