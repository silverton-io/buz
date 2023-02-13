// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package postgresdb

import (
	"testing"

	"github.com/silverton-io/buz/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePostgresDsn(t *testing.T) {
	p := db.ConnectionParams{
		Host: "host",
		Port: 5432,
		Db:   "db",
		User: "usr",
		Pass: "pass",
	}

	expectedDsn := "postgresql://usr:pass@host:5432/db"
	dsn := GenerateDsn(p)

	assert.Equal(t, expectedDsn, dsn)
}
