// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package stats

import "testing"

func TestStatsRoutes(t *testing.T) {
	want := "/stats"
	if STATS_PATH != want {
		t.Fatalf(`Stats path is %v, want %v`, STATS_PATH, want)
	}
}
