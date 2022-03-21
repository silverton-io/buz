package sink

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/tele"
	"github.com/silverton-io/honeypot/pkg/util"
)

var (
	Info = Teal
	Warn = Yellow
	Fata = Red
)

var (
	Black   = Colorize("\033[1;30m%s\033[0m")
	Red     = Colorize("\033[1;31m%s\033[0m")
	Green   = Colorize("\033[1;32m%s\033[0m")
	Yellow  = Colorize("\033[1;33m%s\033[0m")
	Purple  = Colorize("\033[1;34m%s\033[0m")
	Magenta = Colorize("\033[1;35m%s\033[0m")
	Teal    = Colorize("\033[1;36m%s\033[0m")
	White   = Colorize("\033[1;37m%s\033[0m")
)

func Colorize(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type StdoutSink struct{}

func (s *StdoutSink) Initialize(conf config.Sink) {
	log.Debug().Msg("initializing stdout sink")
}

func (s *StdoutSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) {
	if len(validEnvelopes) > 0 {
		validEnvelopes := util.Stringify(validEnvelopes)
		fmt.Println(Green(validEnvelopes))
	}
}

func (s *StdoutSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) {
	if len(invalidEnvelopes) > 0 {
		invalidEnvelopes := util.Stringify(invalidEnvelopes)
		fmt.Println(Red(invalidEnvelopes))
	}
}

func (s *StdoutSink) BatchPublishValidAndInvalid(ctx context.Context, validEnvelopes []envelope.Envelope, invalidEnvelopes []envelope.Envelope, meta *tele.Meta) {
	s.BatchPublishValid(ctx, validEnvelopes)
	s.BatchPublishInvalid(ctx, invalidEnvelopes)
}

func (s *StdoutSink) Close() {
	log.Debug().Msg("closing stdout sink")
}
