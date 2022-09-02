package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateClickhouseDsn(t *testing.T) {
	p := ConnectionParams{
		Host: "myhost",
		Port: 9000,
		Db:   "db",
		User: "usr",
		Pass: "pass",
	}

	dsn := GenerateClickhouseDsn(p)
	expectedDsn := "tcp://myhost:9000?database=db&username=usr&password=pass"

	assert.Equal(t, expectedDsn, dsn)
}
