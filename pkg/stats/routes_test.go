package stats

import "testing"

func TestStatsRoutes(t *testing.T) {
	want := "/stats"
	if STATS_PATH != want {
		t.Fatalf(`Stats path is %v, want %v`, STATS_PATH, want)
	}
}
