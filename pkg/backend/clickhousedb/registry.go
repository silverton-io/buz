// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package clickhousedb

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type RegistryBackend struct {
	gormDb        *gorm.DB
	registryTable string
}

type RegistryTable struct {
	db.BasePKeylessModel
	Name     string `json:"name" gorm:"index:idx_name"`
	Contents string `json:"contents"`
}

func (b *RegistryBackend) Initialize(conf config.Backend) error {
	connParams := db.ConnectionParams{
		Host: conf.ClickhouseHost,
		Port: conf.ClickhousePort,
		Db:   conf.ClickhouseDbName,
		User: conf.ClickhouseUser,
		Pass: conf.ClickhousePass,
	}
	connString := generateDsn(connParams)
	gormDb, err := gorm.Open(clickhouse.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open clickhouse connection")
		return err
	}
	b.gormDb, b.registryTable = gormDb, conf.RegistryTable
	ensureErr := db.EnsureTable(b.gormDb, b.registryTable, RegistryTable{})
	return ensureErr
}

func (b *RegistryBackend) GetRemote(schema string) (contents []byte, err error) {
	var s db.RegistryTable
	b.gormDb.Table(b.registryTable).Where("name = ?", schema).First(&s)
	err = b.gormDb.Error
	if err != nil {
		return nil, err
	}
	return s.Contents, nil
}

func (b *RegistryBackend) Close() {
	log.Info().Msg("🟢 closing clickhouse schema cache backend")
}
