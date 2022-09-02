package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternalRoutes(t *testing.T) {
	assert.Equal(t, STATS_PATH, "/stats")
	assert.Equal(t, HEALTH_PATH, "/health")
	assert.Equal(t, ROUTE_OVERVIEW_PATH, "/routes")
	assert.Equal(t, CONFIG_OVERVIEW_PATH, "/config")
}

func TestSnowplowRoutes(t *testing.T) {
	assert.Equal(t, SNOWPLOW_STANDARD_GET_PATH, "/i")
	assert.Equal(t, SNOWPLOW_STANDARD_POST_PATH, "/com.snowplowanalytics.snowplow/tp2")
	assert.Equal(t, SNOWPLOW_STANDARD_REDIRECT_PATH, "/r/tp2")
}
