package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMysqlDsn(t *testing.T) {
	p := ConnectionParams{
		Host: "host",
		Port: 3306,
		Db:   "db",
		User: "usr",
		Pass: "pass",
	}

	expectedDsn := "usr:pass@tcp(host:3306)/db?parseTime=true"
	dsn := GenerateMysqlDsn(p)

	assert.Equal(t, expectedDsn, dsn)
}
