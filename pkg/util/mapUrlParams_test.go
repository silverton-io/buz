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

func TestMapParams(t *testing.T) {
	url := "something/else?p1=v1&p2=v2&p2=v3"
	want := map[string]interface{}{
		"p1": "v1",
		"p2": "v2",
	}
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = req

	params := MapUrlParams(c)
	assert.Equal(t, params, want)
}

func TestQueryToMap(t *testing.T) {
	params := map[string][]string{
		"p1": {"v1", "v2"},
		"p2": {"v2"},
	}
	expected := map[string]interface{}{
		"p1": "v1",
		"p2": "v2",
	}
	mapped := QueryToMap(params)
	assert.Equal(t, expected, mapped)
}
