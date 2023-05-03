// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/stretchr/testify/assert"
)

var FAKE_IDENTITY_CONF = config.Identity{Cookie: config.IdentityCookie{
	Enabled:  true,
	Name:     "nuid",
	Secure:   true,
	TtlDays:  365,
	Domain:   "new.domain.dev",
	Path:     "/",
	SameSite: "false",
},
	Fallback: "36acc892-3850-4991-bcb6-f8ca7ef23543",
}

var FAKE_RECORDER = httptest.NewRecorder()
var FAKE_CONTEXT, _ = gin.CreateTestContext(FAKE_RECORDER)

func TestGetIdentityOrFallbackNoIdentity(t *testing.T) {
	want := FAKE_IDENTITY_CONF.Fallback
	got := GetIdentityOrFallback(FAKE_CONTEXT, FAKE_IDENTITY_CONF)

	assert.Equal(t, want, got)
}

func TestGetIdentityOrFallbackWithIdentity(t *testing.T) {
	want := "9c96bc72-7913-4058-91bb-3970714a55c2"
	FAKE_CONTEXT.Set(constants.IDENTITY, want)
	got := GetIdentityOrFallback(FAKE_CONTEXT, FAKE_IDENTITY_CONF)

	assert.Equal(t, want, got)
}
