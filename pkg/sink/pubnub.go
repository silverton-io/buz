package sink

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type PubnubSink struct{}

func (s *PubnubSink) Initialize(conf config.Sink) {
	log.Debug().Msg("initializing pubnub sink")
}

func (s *PubnubSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {}

func (s *PubnubSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
}

func (s *PubnubSink) Close() {
	log.Debug().Msg("closing pubnub sink")
}
