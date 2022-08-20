// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package cache

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
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
		log.Error().Err(err).Msg("🔴 could not open pg connection")
		return err
	}
	b.gormDb, b.registryTable = gormDb, conf.RegistryTable
	ensureErr := db.EnsureTable(b.gormDb, b.registryTable, RegistryTable{})
	return ensureErr
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
	log.Info().Msg("🟢 closing postgres schema cache backend")
}
