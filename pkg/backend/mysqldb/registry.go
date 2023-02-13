// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package mysqldb

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type RegistryBackend struct {
	gormDb        *gorm.DB
	registryTable string
}

func (b *RegistryBackend) Initialize(conf config.Backend) error {
	connParams := db.ConnectionParams{
		Host: conf.MysqlHost,
		Port: conf.MysqlPort,
		Db:   conf.MysqlDbName,
		User: conf.MysqlUser,
		Pass: conf.MysqlPass,
	}
	connString := generateDsn(connParams)
	gormDb, err := gorm.Open(mysql.Open(connString))
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open mysql connection")
		return err
	}
	b.gormDb, b.registryTable = gormDb, conf.RegistryTable
	ensureErr := db.EnsureTable(b.gormDb, b.registryTable, db.RegistryTable{})
	return ensureErr
}

func (b *RegistryBackend) GetRemote(schema string) (contents []byte, err error) {
	var s db.RegistryTable
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

func (b *RegistryBackend) Close() {
	log.Info().Msg("🟢 closing mysql schema cache backend")
}
