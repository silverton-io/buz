package util

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestGetSchemaNameFromRequest(t *testing.T) {
	someSchema := "some/schema/v1.0"
	fallback := "some/fallback/v1.0.json"
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	schema := GetSchemaNameFromRequest(c, fallback)

	assert.Equal(t, fallback, schema)

	c.Params = append(c.Params, gin.Param{Key: constants.BUZ_SCHEMA_PARAM, Value: "/" + someSchema})

	schema = GetSchemaNameFromRequest(c, fallback)

	assert.Equal(t, someSchema+JSON_EXTENSION, schema)
}
