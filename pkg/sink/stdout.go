package sink

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
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

func (s *StdoutSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta) {
	if len(validEvents) > 0 {
		validEvents := util.Stringify(validEvents)
		fmt.Println(Green(validEvents))
	}
	if len(invalidEvents) > 0 {
		invalidEvents := util.Stringify(invalidEvents)
		fmt.Println(Red(invalidEvents))
	}
	incrementStats(inputType, len(validEvents), len(invalidEvents), meta)
}

func (s *StdoutSink) Close() {
	log.Debug().Msg("closing stdout sink")
}
