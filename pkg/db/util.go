// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package db

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
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

func EnsureMaterializeWebhookSource(gormDb *gorm.DB, sourceName string, apiKey string, cluster string, model interface{}) error {

	query, err := _generateJsonbQuery(fmt.Sprintf("%s_source", sourceName), "body", model)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not generate view query for " + sourceName)
		return err
	}

	createStatement := fmt.Sprintf(`
		CREATE SECRET IF NOT EXISTS %s_secret AS '%s';
		CREATE SOURCE %s IN CLUSTER %s_source FROM WEBHOOK
		BODY FORMAT JSON
		CHECK (
			WITH (
			  HEADERS, BODY AS request_body,
			  SECRET %s
			)
			constant_time_eq(
			    decode(headers->'x-signature', 'base64'),
			    hmac(request_body, %s_secret, 'sha256')
			)
		);
		CREATE VIEW IF NOT EXISTS %s AS %s
	`, sourceName, apiKey, sourceName, cluster, sourceName, sourceName, sourceName, query)

	err = gormDb.Exec(createStatement).Error

	if err != nil {
		log.Error().Err(err).Msg("🔴 could not create " + sourceName + " source")
		return err
	}
	return nil

}

// Just a guess
func _generateJsonbQuery(tableName, jsonColumn string, structInstance interface{}) (string, error) {
	var selectColumns []string

	val := reflect.ValueOf(structInstance)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		gormTag := field.Tag.Get("gorm")

		if jsonTag == "" {
			continue
		}

		// Handle omitempty for JSONB fields
		if strings.Contains(jsonTag, ",omitempty") {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// Handle multiple options in Gorm tag
		gormOptions := strings.Split(gormTag, ";")
		gormColumnName := jsonTag
		gormType := ""

		for _, option := range gormOptions {
			parts := strings.Split(option, ":")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "column":
					gormColumnName = value
				case "type":
					gormType = value
				}
			}
		}

		if gormType == "" {
			gormType = getGormType(field.Type)
		}
		if gormType == "json" {
			gormType = "jsonb"
		}

		selectColumns = append(selectColumns, fmt.Sprintf(`"%s"."%s"->>'%s'::%s AS "%s"`, tableName, jsonColumn, jsonTag, gormType, gormColumnName))
	}
	if len(selectColumns) == 0 {
		err := fmt.Errorf(`Could not generate query for "%s"`, tableName)
		log.Error().Err(err).Msg("🔴 could not create " + tableName + " source")
		return "", err
	}

	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(selectColumns, ", "), tableName), nil
}

func getGormType(fieldType reflect.Type) string {
	switch fieldType.Kind() {
	case reflect.String:
		return "text"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "numeric"
	case reflect.Bool:
		return "boolean"
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			return "timestamptz"
		}
	}

	return "text"
}
