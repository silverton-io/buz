package cache

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Schema struct {
	gorm.Model
	Name     string         `json:"name" gorm:"index:idx_name,unique"`
	Contents datatypes.JSON `json:"contents"`
}
