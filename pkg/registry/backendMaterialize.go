// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package registry

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
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
		log.Error().Err(err).Msg("🔴 could not open materialize connection")
		return err
	}
	b.gormDb, b.registryTable = gormDb, conf.RegistryTable
	ensureErr := db.EnsureTable(b.gormDb, b.registryTable, RegistryTable{})
	return ensureErr
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
	log.Info().Msg("🟢 closing materialize schema cache backend")
}
