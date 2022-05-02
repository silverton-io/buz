package cache

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/db"
	"github.com/silverton-io/honeypot/pkg/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlSchemaCacheBackend struct {
	gormDb      *gorm.DB
	schemaTable string
}

func (b *MysqlSchemaCacheBackend) Initialize(conf config.Backend) error {
	connParams := db.ConnectionParams{
		Host: conf.MysqlHost,
		Port: conf.MysqlPort,
		Db:   conf.MysqlDbName,
		User: conf.MysqlUser,
		Pass: conf.MysqlPass,
	}
	connString := db.GenerateMysqlDsn(connParams)
	gormDb, err := gorm.Open(mysql.Open(connString))
	if err != nil {
		log.Error().Err(err).Msg("could not open mysql connection")
		return err
	}
	b.gormDb, b.schemaTable = gormDb, conf.SchemaTable
	schemaTblExists := b.gormDb.Migrator().HasTable(b.schemaTable)
	if !schemaTblExists {
		log.Debug().Msg(b.schemaTable + " table doesn't exist - creating")
		err = b.gormDb.Table(b.schemaTable).AutoMigrate(Schema{})
		if err != nil {
			log.Error().Err(err).Msg("could not create schema table")
			return err
		}
	}
	return nil
}

func (b *MysqlSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	var s Schema
	b.gormDb.Table(b.schemaTable).Where("name = ?", schema).First(&s)
	err = b.gormDb.Error
	if err != nil {
		log.Error().Err(err).Msg("gorm error")
		return nil, err
	}
	util.Pprint(s)
	contents, err = json.Marshal(s.Contents)
	if err != nil {
		log.Error().Err(err).Msg("could not marshal schema contents")
	}
	return contents, nil
}

func (b *MysqlSchemaCacheBackend) Close() {
	log.Info().Msg("closing mysql schema cache backend")
}
