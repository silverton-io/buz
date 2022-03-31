package snowplow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	assert.Equal(t, DEFAULT_GET_PATH, "i")
	assert.Equal(t, DEFAULT_POST_PATH, "com.snowplowanalytics.snowplow/tp2")
	assert.Equal(t, DEFAULT_REDIRECT_PATH, "r/tp2")
}
