package health

import "testing"

func TestHealthPath(t *testing.T) {
	want := "/health"
	if HEALTH_PATH != want {
		t.Fatalf(`Health path is %v, want %v`, HEALTH_PATH, want)
	}
}
