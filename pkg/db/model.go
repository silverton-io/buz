// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package db

import (
	"time"

	"gorm.io/datatypes"
)

type BasePKeylessModel struct {
	CreatedAt time.Time  `json:"-" sql:"index"`
	UpdatedAt time.Time  `json:"-" sql:"index"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

type RegistryTable struct {
	BasePKeylessModel
	Name     string         `json:"name" gorm:"index:idx_name"`
	Contents datatypes.JSON `json:"contents"`
}
