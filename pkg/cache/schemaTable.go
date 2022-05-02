package cache

import (
	"github.com/jackc/pgtype"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Schema struct {
	gorm.Model
	Name     string         `json:"name" gorm:"index:idx_name,unique"`
	Contents datatypes.JSON `json:"contents" gorm:"type:json"`
}

type PgSchema struct {
	gorm.Model
	Name     string       `json:"name" gorm:"index:idx_name,unique"`
	Contents pgtype.JSONB `json:"contents" gorm:"type:jsonb"` // Again, hate to do it this way when the only difference is gorm "json" vs "jsonb" tags...
}
