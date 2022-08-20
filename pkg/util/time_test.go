// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"testing"
	"time"
)

func TestGetDuration(t *testing.T) {
	want := 2 * time.Second
	start := time.Now().UTC()
	end := start.Add(want)
	duration := GetDuration(start, end)
	if duration != want {
		t.Fatalf(`GetDuration(%v, %v) = %v, want %v`, start, end, duration, want)
	}
}
