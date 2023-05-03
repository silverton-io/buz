// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/response"
)

func TestTimeout(t *testing.T) {
	u := "/somepath"
	fastTimeout := config.Timeout{
		Enabled: true,
		Ms:      1,
	}
	slowTimeout := config.Timeout{
		Enabled: true,
		Ms:      30,
	}
	okResponse, _ := json.Marshal(response.Ok)
	timeoutResponse, _ := json.Marshal(response.Timeout)

	var testCases = []struct {
		config       config.Timeout
		wantCode     int
		wantResponse []byte
	}{
		{config: fastTimeout, wantCode: 408, wantResponse: timeoutResponse},
		{config: slowTimeout, wantCode: 200, wantResponse: okResponse},
	}

	for _, tc := range testCases {
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(Timeout(tc.config))
		r.GET(u, testHandler)
		ts := httptest.NewServer(r)

		resp, _ := http.Get(ts.URL + u)
		if resp.StatusCode != tc.wantCode {
			t.Fatalf(`got status code %v, want %v`, resp.StatusCode, tc.wantCode)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		bodyEquiv := reflect.DeepEqual(body, tc.wantResponse)
		if !bodyEquiv {
			t.Fatalf(`got response %v, want %v`, body, tc.wantResponse)
		}
	}
}
