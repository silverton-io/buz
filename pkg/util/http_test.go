// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHttpHeadersToMap ensures
// only the first value of a specific header
// shows up.
func TestHttpHeadersToMap(t *testing.T) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	want := map[string]interface{}{
		"Header1": "val1",
		"Header2": []string{"val2", "val3"},
	}
	c.Request.Header.Add("header1", "val1")
	c.Request.Header.Add("Header2", "val2")
	c.Request.Header.Add("Header2", "val3")

	got := HttpHeadersToMap(c)

	assert.Equal(t, want, got)
}
