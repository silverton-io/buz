// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCors(t *testing.T) {
	u := "/test"
	conf := config.Cors{
		Enabled:          true,
		AllowOrigin:      []string{"http://allowed-origin.com"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "OPTIONS"},
		MaxAge:           86400,
	}
	r := gin.New()
	r.Use(CORS(conf))
	r.GET(u, testHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("preflight success", func(t *testing.T) {
		var client = &http.Client{}
		req, _ := http.NewRequest(http.MethodOptions, ts.URL+u, nil)
		req.Header.Set("Origin", "http://allowed-origin.com")
		resp, _ := client.Do(req)

		assert.Equal(t, []string{"true"}, resp.Header["Access-Control-Allow-Credentials"])
		assert.Equal(t, []string{"GET, OPTIONS"}, resp.Header["Access-Control-Allow-Methods"])
		assert.Equal(t, []string{"http://allowed-origin.com"}, resp.Header["Access-Control-Allow-Origin"])
		assert.Equal(t, []string{"86400"}, resp.Header["Access-Control-Max-Age"])
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("preflight fail", func(t *testing.T) {
		var client = &http.Client{}
		req, _ := http.NewRequest(http.MethodOptions, ts.URL+u, nil)
		req.Header.Set("Origin", "http://not-allowed-origin.com")
		resp, _ := client.Do(req)

		assert.Equal(t, []string{"true"}, resp.Header["Access-Control-Allow-Credentials"])
		assert.Equal(t, []string{"GET, OPTIONS"}, resp.Header["Access-Control-Allow-Methods"])
		assert.Equal(t, []string([]string(nil)), resp.Header["Access-Control-Allow-Origin"])
		assert.Equal(t, []string{"86400"}, resp.Header["Access-Control-Max-Age"])
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("get", func(t *testing.T) {
		resp, _ := http.Get(ts.URL + u)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
