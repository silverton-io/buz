package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"

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
	cfg := elasticsearch.Config{
		Addresses: conf.ElasticsearchHosts,
		Username:  conf.ElasticsearchUsername,
		Password:  conf.ElasticsearchPassword,
	}
	es, _ := elasticsearch.NewClient(cfg)
	s.client, s.validIndex, s.invalidIndex = es, conf.ValidIndex, conf.InvalidIndex
}

func (s *ElasticsearchSink) batchPublish(ctx context.Context, index string, envelopes []envelope.Envelope) {
	var wg sync.WaitGroup
	for _, envelope := range envelopes {
		eByte, err := json.Marshal(envelope)
		reader := bytes.NewReader(eByte)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not encode envelope to buffer")
		} else {
			wg.Add(1)
			envId := envelope.Id.String()
			_, err := s.client.Create(index, envId, reader)
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not publish envelope to elasticsearch: " + envId)
			} else {
				log.Debug().Msg("published envelope to " + index + " index: " + envId)
			}
			defer wg.Done()
		}
	}
	wg.Wait()
}

func (s *ElasticsearchSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.validIndex, validEnvelopes)
}

func (s *ElasticsearchSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.invalidIndex, invalidEnvelopes)
}

func (s *ElasticsearchSink) BatchPublishValidAndInvalid(ctx context.Context, validEnvelopes []envelope.Envelope, invalidEnvelopes []envelope.Envelope, meta *tele.Meta) {
	go s.BatchPublishValid(ctx, validEnvelopes)
	go s.BatchPublishInvalid(ctx, invalidEnvelopes)
	// FIXME!! Write envelope publish stats
}

func (s *ElasticsearchSink) Close() {
	log.Debug().Msg("closing elasticsearch sink client")
}
