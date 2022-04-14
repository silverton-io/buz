package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type ElasticsearchSink struct {
	id           *uuid.UUID
	name         string
	client       *elasticsearch.Client
	validIndex   string
	invalidIndex string
}

func (s *ElasticsearchSink) Id() *uuid.UUID {
	return s.id
}

func (s *ElasticsearchSink) Name() string {
	return s.name
}

func (s *ElasticsearchSink) Initialize(conf config.Sink) error {
	cfg := elasticsearch.Config{
		Addresses: conf.ElasticsearchHosts,
		Username:  conf.ElasticsearchUsername,
		Password:  conf.ElasticsearchPassword,
	}
	es, err := elasticsearch.NewClient(cfg)
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	s.client, s.validIndex, s.invalidIndex = es, conf.ValidIndex, conf.InvalidIndex
	return err
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
			envId := envelope.Uuid.String()
			_, err := s.client.Create(index, envId, reader)
			if err != nil {
				log.Error().Stack().Interface("envelopeId", envId).Err(err).Msg("could not publish envelope to elasticsearch")
			} else {
				log.Debug().Interface("envelopeId", envId).Interface("indexId", index).Msg("published envelope to index")
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

func (s *ElasticsearchSink) Close() {
	log.Debug().Msg("closing elasticsearch sink client")
}
