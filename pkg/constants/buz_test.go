// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuzConstants(t *testing.T) {
	assert.Equal(t, "hps", BUZ_SCHEMA_PARAM)
	assert.Equal(t, "hpbp", BUZ_BASE64_ENCODED_PAYLOAD_PARAM)
	assert.Equal(t, "unknown", UNKNOWN)
}
