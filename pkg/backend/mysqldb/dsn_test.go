// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package mysqldb

import (
	"testing"

	"github.com/silverton-io/buz/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDsn(t *testing.T) {
	p := db.ConnectionParams{
		Host: "host",
		Port: 3306,
		Db:   "db",
		User: "usr",
		Pass: "pass",
	}

	expectedDsn := "usr:pass@tcp(host:3306)/db?parseTime=true"
	dsn := generateDsn(p)

	assert.Equal(t, expectedDsn, dsn)
}
