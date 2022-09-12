// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

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
	assert.Equal(t, "pixel", PIXEL)
}
