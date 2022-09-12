package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbConstants(t *testing.T) {
	assert.Equal(t, "postgres", POSTGRES)
	assert.Equal(t, "materialize", MATERIALIZE)
	assert.Equal(t, "mysql", MYSQL)
	assert.Equal(t, "clickhouse", CLICKHOUSE)
	assert.Equal(t, "mongodb", MONGODB)
	assert.Equal(t, "elasticsearch", ELASTICSEARCH)
	assert.Equal(t, "timescale", TIMESCALE)
}
