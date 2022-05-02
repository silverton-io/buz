package cache

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresSchemaCacheBackend struct {
	gormDb      *gorm.DB
	schemaTable string
}

func (b *PostgresSchemaCacheBackend) Initialize(conf config.Backend) error {
	connParams := db.ConnectionParams{
		Host: conf.PgHost,
		Port: conf.PgPort,
		Db:   conf.PgDbName,
		User: conf.PgUser,
		Pass: conf.PgPass,
	}
	connString := db.GeneratePostgresDsn(connParams)
	gormDb, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("could not open pg connection")
		return err
	}
	b.gormDb, b.schemaTable = gormDb, conf.SchemaTable
	schemaTblExists := b.gormDb.Migrator().HasTable(b.schemaTable)
	if !schemaTblExists {
		log.Debug().Msg(b.schemaTable + " table doesn't exist - creating")
		err = b.gormDb.Table(b.schemaTable).AutoMigrate(PgSchema{})
		if err != nil {
			log.Error().Err(err).Msg("could not create schema table")
			return err
		}
	}
	return nil
}

func (b *PostgresSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	var s PgSchema
	b.gormDb.Table(b.schemaTable).Where("name = ?", schema).First(&s)
	err = b.gormDb.Error
	if err != nil {
		return nil, err
	}
	return s.Contents.Bytes, nil
}

func (b *PostgresSchemaCacheBackend) Close() {
	log.Info().Msg("closing postgres schema cache backend")
}
