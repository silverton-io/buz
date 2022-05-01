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
	connParams := db.DbConnectionParams{
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
		}
	}
	return nil
}

func (b *PostgresSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	var s PgSchema
	res := b.gormDb.Where("name = ?", schema).First(&s)
}

func (b *PostgresSchemaCacheBackend) Close() {

}
