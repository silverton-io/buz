// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package db

import (
	// "github.com/silverton-io/buz/pkg/db"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// database tables with string columns.
type StringEnvelope struct {
	Uuid            uuid.UUID         `json:"uuid" gorm:"type:uuid"`
	Timestamp       time.Time         `json:"timestamp" sql:"index"`
	BuzVersion      string            `json:"buzVersion"`
	IsValid         bool              `json:"isValid"`
	ValidationError map[string]string `json:"validationError,omitempty" gorm:"type:json"`
}

func Test_generateJsonbQeury(t *testing.T) {
	q, err := _generateJsonbQuery("foo", "bar", StringEnvelope{})
	assert.Nil(t, err, "Error should be nil when generating jsonb query")
	test_response := `SELECT "foo"."bar"->>'uuid'::uuid AS "uuid", "foo"."bar"->>'timestamp'::timestamptz AS "timestamp", "foo"."bar"->>'buzVersion'::text AS "buzVersion", "foo"."bar"->>'isValid'::boolean AS "isValid", "foo"."bar"->>'validationError'::jsonb AS "validationError" FROM foo`
	assert.Equal(t, q, test_response)
}
