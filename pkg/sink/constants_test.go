// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSinkConstants(t *testing.T) {
	assert.Equal(t, "buz_valid_events", BUZ_VALID_EVENTS)
	assert.Equal(t, "buz_invalid_events", BUZ_INVALID_EVENTS)
}
