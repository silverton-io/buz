package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternalRoutes(t *testing.T) {
	assert.Equal(t, "/stats", STATS_PATH)
	assert.Equal(t, "/health", HEALTH_PATH)
	assert.Equal(t, "/routes", ROUTE_OVERVIEW_PATH)
	assert.Equal(t, "/config", CONFIG_OVERVIEW_PATH)
}

func TestSnowplowRoutes(t *testing.T) {
	assert.Equal(t, "/i", SNOWPLOW_STANDARD_GET_PATH)
	assert.Equal(t, "/com.snowplowanalytics.snowplow/tp2", SNOWPLOW_STANDARD_POST_PATH)
	assert.Equal(t, "/r/tp2", SNOWPLOW_STANDARD_REDIRECT_PATH)
}
