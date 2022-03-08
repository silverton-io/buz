package sink

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/tele"
	"github.com/silverton-io/honeypot/pkg/util"
)

type StdoutSink struct{}

func (s *StdoutSink) Initialize(conf config.Sink) {
	log.Debug().Msg("initializing stdout sink")
}

func (s *StdoutSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta) {
	fmt.Println("Valid events:")
	util.PrettyPrint(validEvents)
	fmt.Println("Invalid events:")
	util.PrettyPrint(invalidEvents)
	incrementStats(inputType, len(validEvents), len(invalidEvents), meta)
}

func (s *StdoutSink) Close() {
	log.Debug().Msg("closing stdout sink")
}
