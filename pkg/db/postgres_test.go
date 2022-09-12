package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePostgresDsn(t *testing.T) {
	p := ConnectionParams{
		Host: "host",
		Port: 5432,
		Db:   "db",
		User: "usr",
		Pass: "pass",
	}

	expectedDsn := "postgresql://usr:pass@host:5432/db"
	dsn := GeneratePostgresDsn(p)

	assert.Equal(t, expectedDsn, dsn)
}
