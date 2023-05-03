// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/meta"
	testutil "github.com/silverton-io/buz/pkg/testUtil"
	"github.com/stretchr/testify/assert"
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

	rec, c, _ := testutil.BuildRecordedEngine()

	handler := StatsHandler(&m)

	handler(c)

	resp := rec.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
	expectedResponse := StatsResponse{
		CollectorMeta: &m,
		// Stats:         &s,
	}
	expected, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(`Could not marshal expected response`)
	}
	assert.Equal(t, expected, b)
}
