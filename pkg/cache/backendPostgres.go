package cache

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresSchemaCacheBackend struct {
	gormDb        *gorm.DB
	registryTable string
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
	b.gormDb, b.registryTable = gormDb, conf.RegistryTable
	schemaTblExists := b.gormDb.Migrator().HasTable(b.registryTable)
	if !schemaTblExists {
		log.Debug().Msg(b.registryTable + " table doesn't exist - creating")
		err = b.gormDb.Table(b.registryTable).AutoMigrate(RegistryTable{})
		if err != nil {
			log.Error().Err(err).Msg("could not create schema table")
			return err
		}
	}
	return nil
}

func (b *PostgresSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	var s RegistryTable
	b.gormDb.Table(b.registryTable).Where("name = ?", schema).First(&s)
	err = b.gormDb.Error
	if err != nil {
		return nil, err
	}
	return s.Contents, nil
}

func (b *PostgresSchemaCacheBackend) Close() {
	log.Info().Msg("closing postgres schema cache backend")
}
