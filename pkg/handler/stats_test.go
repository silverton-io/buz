// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/stats"
)

func TestStatsHandler(t *testing.T) {
	u := uuid.New()
	now := time.Now().UTC()
	m := meta.CollectorMeta{
		Version:       "1.0.x",
		InstanceId:    u,
		StartTime:     now,
		TrackerDomain: "somewhere.net",
		CookieDomain:  "somewhere.io",
	}
	s := stats.ProtocolStats{}
	s.Build()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	handler := StatsHandler(&m, &s)

	handler(c)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf(`StatsHandler returned %d, want %d`, resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
	expectedResponse := StatsResponse{
		CollectorMeta: &m,
		Stats:         &s,
	}
	expected, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(`Could not marshal expected response`)
	}
	equiv := reflect.DeepEqual(b, expected)
	if !equiv {
		t.Fatalf(`StatsHandler returned %v, want %v`, b, expected)
	}
}
