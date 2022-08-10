// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCors(t *testing.T) {
	u := "/test"
	conf := config.Cors{
		Enabled:          true,
		AllowOrigin:      []string{"*"},
		AllowCredentials: false,
		AllowMethods:     []string{"GET", "OPTIONS"},
		MaxAge:           86400,
	}
	r := gin.New()
	r.Use(CORS(conf))
	r.GET(u, testHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("preflight", func(t *testing.T) {
		var client = &http.Client{}
		req, _ := http.NewRequest(http.MethodOptions, ts.URL+u, nil)
		resp, _ := client.Do(req)

		assert.Equal(t, []string{"false"}, resp.Header["Access-Control-Allow-Credentials"])
		assert.Equal(t, []string{"GET, OPTIONS"}, resp.Header["Access-Control-Allow-Methods"])
		assert.Equal(t, []string{"*"}, resp.Header["Access-Control-Allow-Origin"])
		assert.Equal(t, []string{"86400"}, resp.Header["Access-Control-Max-Age"])
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("get", func(t *testing.T) {
		resp, _ := http.Get(ts.URL + u)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
