package stats

import "testing"

func TestStatsRoutes(t *testing.T) {
	want := "/stats"
	if STATS_PATH != want {
		t.Fatalf(`Got %v, want %v`, STATS_PATH, want)
	}
}
