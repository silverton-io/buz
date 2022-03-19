package sink

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/tele"
)

type ElasticsearchSink struct {
	client       *elasticsearch.Client
	validIndex   string
	invalidIndex string
}

func (s *ElasticsearchSink) Initialize(conf config.Sink) {
	es, _ := elasticsearch.NewDefaultClient()
	s.client = es
}

func (s *ElasticsearchSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {}

func (s *ElasticsearchSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {}

func (s *ElasticsearchSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEnvelopes []envelope.Envelope, invalidEnvelopes []envelope.Envelope, meta *tele.Meta) {
}

func (s *ElasticsearchSink) Close() {
	log.Debug().Msg("closing elasticsearch sink client")
}
