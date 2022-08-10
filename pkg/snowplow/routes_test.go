// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

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
