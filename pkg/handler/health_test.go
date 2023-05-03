// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/silverton-io/buz/pkg/response"
	testutil "github.com/silverton-io/buz/pkg/testUtil"
)

func TestHealthcheckHandler(t *testing.T) {
	rec, c, _ := testutil.BuildRecordedEngine()

	HealthcheckHandler(c)

	resp := rec.Result()

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf(`HealthcheckHandler returned status code %v, want %v`, resp.StatusCode, http.StatusOK)
	}
	b, _ := io.ReadAll(resp.Body)
	marshaledB, _ := json.Marshal(response.Ok)
	equiv := reflect.DeepEqual(b, marshaledB)
	if !equiv {
		t.Fatalf(`HealthcheckHandler returned body %v, want %v`, b, marshaledB)
	}
}
