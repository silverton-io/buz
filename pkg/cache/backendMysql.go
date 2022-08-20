// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package cache

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlSchemaCacheBackend struct {
	gormDb        *gorm.DB
	registryTable string
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
		log.Error().Err(err).Msg("🔴 could not open mysql connection")
		return err
	}
	b.gormDb, b.registryTable = gormDb, conf.RegistryTable
	ensureErr := db.EnsureTable(b.gormDb, b.registryTable, RegistryTable{})
	return ensureErr
}

func (b *MysqlSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	var s RegistryTable
	b.gormDb.Table(b.registryTable).Where("name = ?", schema).First(&s)
	err = b.gormDb.Error
	if err != nil {
		log.Error().Err(err).Msg("🔴 gorm error")
		return nil, err
	}
	contents, err = json.Marshal(s.Contents)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not marshal schema contents")
	}
	return contents, nil
}

func (b *MysqlSchemaCacheBackend) Close() {
	log.Info().Msg("🟢 closing mysql schema cache backend")
}
