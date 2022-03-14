package input

import "testing"

func TestInputConst(t *testing.T) {
	want_sp := "snowplow"
	want_gen := "generic"
	want_ce := "cloudevents"

	t.Run("snowplow", func(t *testing.T) {
		if SNOWPLOW_INPUT != want_sp {
			t.Fatalf(`got %v, want %v`, SNOWPLOW_INPUT, want_sp)
		}
	})
	t.Run("generic", func(t *testing.T) {
		if GENERIC_INPUT != want_gen {
			t.Fatalf(`got %v, want %v`, GENERIC_INPUT, want_gen)
		}
	})
	t.Run("cloudevents", func(t *testing.T) {
		if CLOUDEVENTS_INPUT != want_ce {
			t.Fatalf(`got %v, want %v`, CLOUDEVENTS_INPUT, want_ce)
		}
	})
}
