package cache

import (
	"github.com/silverton-io/honeypot/pkg/db"
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
