// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestGetDurationFromString(t *testing.T) {
	var testCases = []struct {
		period       string
		wantDuration time.Duration
	}{
		{"MS", 1 * time.Millisecond},
		{"S", 1 * time.Second},
		{"M", 1 * time.Minute},
		{"H", 1 * time.Hour},
		{"D", 24 * time.Hour},
		{"default", 1 * time.Second},
	}

	for _, tc := range testCases {
		t.Run(tc.period, func(t *testing.T) {
			got := getDurationFromString(tc.period)
			assert.Equal(t, tc.wantDuration, got)
		})
	}
}

func TestOnLimitReachedHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	onLimitReachedHandler(c)

	body, _ := ioutil.ReadAll(rec.Body)
	wantBody, _ := json.Marshal(response.RateLimitExceeded)
	assert.Equal(t, http.StatusTooManyRequests, rec.Result().StatusCode)
	assert.Equal(t, wantBody, body)
}

func TestBuildRateLimiter(t *testing.T) {
	c := config.RateLimiter{
		Enabled: true,
		Period:  "H",
		Limit:   int64(1),
	}
	wantDuration := getDurationFromString(c.Period)
	limiter := BuildRateLimiter(c)
	assert.Equal(t, limiter.Rate.Period, wantDuration)
	assert.Equal(t, limiter.Rate.Limit, c.Limit)
}

func TestBuildRateLimiterMiddleware(t *testing.T) {
	c := config.RateLimiter{
		Enabled: true,
		Period:  "H",
		Limit:   int64(1),
	}
	limiter := BuildRateLimiter(c)
	BuildRateLimiterMiddleware(limiter)
}
