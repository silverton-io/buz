package protocol

import "testing"

func TestInputConst(t *testing.T) {
	want_sp := "snowplow"
	want_gen := "generic"
	want_ce := "cloudevents"

	t.Run("snowplow", func(t *testing.T) {
		if SNOWPLOW != want_sp {
			t.Fatalf(`got %v, want %v`, SNOWPLOW, want_sp)
		}
	})
	t.Run("generic", func(t *testing.T) {
		if GENERIC != want_gen {
			t.Fatalf(`got %v, want %v`, GENERIC, want_gen)
		}
	})
	t.Run("cloudevents", func(t *testing.T) {
		if CLOUDEVENTS != want_ce {
			t.Fatalf(`got %v, want %v`, CLOUDEVENTS, want_ce)
		}
	})
}
