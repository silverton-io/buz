// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package clickhousedb

import (
	"testing"

	"github.com/silverton-io/buz/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDsn(t *testing.T) {
	p := db.ConnectionParams{
		Host: "myhost",
		Port: 9000,
		Db:   "db",
		User: "usr",
		Pass: "pass",
	}

	dsn := generateDsn(p)
	expectedDsn := "tcp://myhost:9000?database=db&username=usr&password=pass"

	assert.Equal(t, expectedDsn, dsn)
}
