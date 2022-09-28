// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package db

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// EnsureTable creates a table according to the specified model if it
// does not already exist.
func EnsureTable(gormDb *gorm.DB, tableName string, model interface{}) error {
	tblExists := gormDb.Migrator().HasTable(tableName)
	if !tblExists {
		log.Debug().Msg("🟡 " + tableName + " table doesn't exist - creating")
		err := gormDb.Table(tableName).AutoMigrate(model)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not create " + tableName + " table")
		}
	} else {
		log.Debug().Msg("🟡 " + tableName + " table already exists - not creating")
	}
	return nil
}
