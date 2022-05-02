package sink

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/db"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

func generateClickhouseDsn(conf config.Sink) string {
	// "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
	port := strconv.FormatUint(uint64(conf.ClickhousePort), 10)
	return "tcp://" + conf.ClickhouseHost + ":" + port + "?database=" + conf.ClickhouseDbName + "&username=" + conf.ClickhouseUser + "&password=" + conf.ClickhousePass
}

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
	connString := generateClickhouseDsn(conf)
	gormDb, err := gorm.Open(clickhouse.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("could not open clickhouse connection")
		return err
	}
	s.gormDb, s.validTable, s.invalidTable = gormDb, conf.ValidTable, conf.InvalidTable
	for _, tbl := range []string{s.validTable, s.invalidTable} {
		tblExists := s.gormDb.Migrator().HasTable(tbl)
		if !tblExists {
			log.Debug().Msg(tbl + " table doesn't exist - ensuring")
			err = s.gormDb.Table(tbl).AutoMigrate(&envelope.ClickhouseEnvelope{})
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
