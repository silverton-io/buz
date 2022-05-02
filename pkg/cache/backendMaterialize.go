package cache

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MaterializeSchemaCacheBackend struct {
	gormDb        *gorm.DB
	registryTable string
}

func (b *MaterializeSchemaCacheBackend) Initialize(conf config.Backend) error {
	connParams := db.ConnectionParams{
		Host: conf.MzHost,
		Port: conf.MzPort,
		Db:   conf.MzDbName,
		User: conf.MzUser,
		Pass: conf.MzPass,
	}
	connString := db.GenerateMzDsn(connParams)
	gormDb, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("could not open materialize connection")
		return err
	}
	b.gormDb, b.registryTable = gormDb, conf.RegistryTable
	schemaTblExists := b.gormDb.Migrator().HasTable(b.registryTable)
	if !schemaTblExists {
		log.Debug().Msg(b.registryTable + " table doesn't exist - creating")
		err = b.gormDb.Table(b.registryTable).AutoMigrate(RegistryTable{})
		if err != nil {
			log.Error().Err(err).Msg("could not create table")
			return err
		}
	}
	return nil
}

func (b *MaterializeSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	var s RegistryTable
	b.gormDb.Table(b.registryTable).Where("name = ?", schema).First(&s)
	err = b.gormDb.Error
	if err != nil {
		return nil, err
	}
	return s.Contents, nil
}

func (b *MaterializeSchemaCacheBackend) Close() {
	log.Info().Msg("closing materialize schema cache backend")
}
