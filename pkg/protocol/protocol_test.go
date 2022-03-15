package protocol

import "testing"

func TestInputConst(t *testing.T) {
	want_sp := "snowplow"
	want_gen := "generic"
	want_ce := "cloudevents"

	t.Run("snowplow", func(t *testing.T) {
		if SNOWPLOW_PROTOCOL != want_sp {
			t.Fatalf(`got %v, want %v`, SNOWPLOW_PROTOCOL, want_sp)
		}
	})
	t.Run("generic", func(t *testing.T) {
		if GENERIC_PROTOCOL != want_gen {
			t.Fatalf(`got %v, want %v`, GENERIC_PROTOCOL, want_gen)
		}
	})
	t.Run("cloudevents", func(t *testing.T) {
		if CLOUDEVENTS_PROTOCOL != want_ce {
			t.Fatalf(`got %v, want %v`, CLOUDEVENTS_PROTOCOL, want_ce)
		}
	})
}
