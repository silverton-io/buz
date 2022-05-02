package cache

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Schema struct {
	gorm.Model
	Name   string      `json:"name" gorm:"index:idx_name"`
	Schema interface{} `json:"schema" gorm:"type:json"`
}

func (s Schema) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

func (s Schema) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &s)
}

type PgSchema struct {
	gorm.Model
	Name   string       `json:"name" gorm:"index:idx_name,unique"`
	Schema pgtype.JSONB `json:"schema" gorm:"type:jsonb"` // Again, hate to do it this way when the only difference is gorm "json" vs "jsonb" tags...
}
