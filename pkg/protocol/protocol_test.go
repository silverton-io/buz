package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputConst(t *testing.T) {
	assert.Equal(t, "snowplow", SNOWPLOW)
	assert.Equal(t, "generic", GENERIC)
	assert.Equal(t, "cloudevents", CLOUDEVENTS)
	assert.Equal(t, "webhook", WEBHOOK)
}
