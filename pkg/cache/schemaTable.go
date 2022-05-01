package cache

import "gorm.io/gorm"

type Schema struct {
	gorm.Model
	name   string                 `json:"name"`
	schema map[string]interface{} `json:"schema" gorm:"type:json"`
}

type PgSchema struct {
	gorm.Model
	name   string                 `json:"name" gorm:"index"`
	schema map[string]interface{} `json:"schema" gorm:"type:jsonb"`
}
