// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package duckdb

import (
	"time"

	"github.com/alifiroozi80/duckdb"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RegistryBackend struct {
	gormDb        *gorm.DB
	registryTable string
}

// We will create this 'RegistryTable' struct to get rid of `deleted_at`
// Becuase of DuckDB gorm limitation.
// https://github.com/alifiroozi80/duckdb?tab=readme-ov-file#limitations
type RegistryTable struct {
	CreatedAt time.Time      `json:"-" sql:"index"`
	UpdatedAt time.Time      `json:"-" sql:"index"`
	Name      string         `json:"name" gorm:"index:idx_name"`
	Contents  datatypes.JSON `json:"contents"`
}

func (b *RegistryBackend) Initialize(conf config.Backend) error {
	gormDb, err := gorm.Open(duckdb.Open("duckdb.ddb"), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open ddb connection")
		return err
	}
	b.gormDb, b.registryTable = gormDb, conf.RegistryTable
	ensureErr := db.EnsureTable(b.gormDb, b.registryTable, RegistryTable{})
	return ensureErr
}

func (b *RegistryBackend) GetRemote(schema string) (contents []byte, err error) {
	var s RegistryTable
	b.gormDb.Table(b.registryTable).Where("name = ?", schema).First(&s)
	err = b.gormDb.Error
	if err != nil {
		return nil, err
	}
	return s.Contents, nil
}

func (b *RegistryBackend) Close() {
	log.Info().Msg("🟢 closing duckdb schema cache backend")
}
