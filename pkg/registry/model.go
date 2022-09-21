// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package registry

import (
	"github.com/silverton-io/buz/pkg/db"
	"gorm.io/datatypes"
)

type RegistryTable struct {
	db.BasePKeylessModel
	Name     string         `json:"name" gorm:"index:idx_name"`
	Contents datatypes.JSON `json:"contents"`
}

type ClickhouseRegistryTable struct {
	db.BasePKeylessModel
	Name     string `json:"name" gorm:"index:idx_name"`
	Contents string `json:"contents"`
}
