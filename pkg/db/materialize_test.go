package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMzDsn(t *testing.T) {
	p := ConnectionParams{
		Host: "host",
		Port: 6875,
		Db:   "db",
		User: "usr",
		Pass: "pass",
	}

	expectedDsn := "postgresql://usr:pass@host:6875/db"
	dsn := GenerateMzDsn(p)

	assert.Equal(t, expectedDsn, dsn)
}
