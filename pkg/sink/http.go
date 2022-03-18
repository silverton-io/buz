package sink

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/tele"
)

type HttpSink struct{}

func (s *HttpSink) Initialize(conf config.Sink) {
	log.Debug().Msg("initializing http sink")
}

func (s *HttpSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) {}

func (s *HttpSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) {}

func (s *HttpSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEnvelopes []envelope.Envelope, invalidEnvelopes []envelope.Envelope, meta *tele.Meta) {
}

func (s *HttpSink) Close() {}
