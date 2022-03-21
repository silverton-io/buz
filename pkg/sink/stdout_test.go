package sink

import (
	"context"
	"testing"

	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/tele"
	"github.com/stretchr/testify/assert"
)

func TestConst(t *testing.T) {
	var testCases = []struct {
		name string
		have func(...interface{}) string
		want func(...interface{}) string
	}{
		{"black", Black, Colorize("\033[1;30m%s\033[0m")},
		{"red", Red, Colorize("\033[1;31m%s\033[0m")},
		{"green", Green, Colorize("\033[1;32m%s\033[0m")},
		{"yellow", Yellow, Colorize("\033[1;33m%s\033[0m")},
		{"purple", Purple, Colorize("\033[1;34m%s\033[0m")},
		{"magenta", Magenta, Colorize("\033[1;35m%s\033[0m")},
		{"teal", Teal, Colorize("\033[1;36m%s\033[0m")},
		{"white", White, Colorize("\033[1;37m%s\033[0m")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			have := tc.have("hello")
			want := tc.want("hello")
			assert.Equal(t, want, have)
		})
	}
}

func TestColorize(t *testing.T) {
	s := "hello"
	want := "\x1b[1;36mhello\x1b[0m"
	have := Colorize("\033[1;36m%s\033[0m")(s)
	assert.Equal(t, want, have)
}

func TestStdoutSink(t *testing.T) {
	c := config.Sink{
		Type: STDOUT,
	}
	m := tele.Meta{}
	ctx := context.Background()
	sink := StdoutSink{}

	sink.Initialize(c)
	sink.BatchPublishValidAndInvalid(ctx, []envelope.Envelope{}, []envelope.Envelope{}, &m)
	sink.Close()
}
