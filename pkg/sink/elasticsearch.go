// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
	"github.com/silverton-io/buz/pkg/envelope"
)

type ElasticsearchSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	client           *elasticsearch.Client
	validIndex       string
	invalidIndex     string
}

func (s *ElasticsearchSink) Id() *uuid.UUID {
	return s.id
}

func (s *ElasticsearchSink) Name() string {
	return s.name
}

func (s *ElasticsearchSink) Type() string {
	return db.ELASTICSEARCH
}

func (s *ElasticsearchSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *ElasticsearchSink) Initialize(conf config.Sink) error {
	cfg := elasticsearch.Config{
		Addresses: conf.ElasticsearchHosts,
		Username:  conf.ElasticsearchUsername,
		Password:  conf.ElasticsearchPassword,
	}
	es, err := elasticsearch.NewClient(cfg)
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.client, s.validIndex, s.invalidIndex = es, conf.ValidIndex, conf.InvalidIndex
	return err
}

func (s *ElasticsearchSink) batchPublish(ctx context.Context, index string, envelopes []envelope.Envelope) error {
	var wg sync.WaitGroup
	for _, envelope := range envelopes {
		eByte, err := json.Marshal(envelope)
		reader := bytes.NewReader(eByte)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not encode envelope to buffer")
			return err
		} else {
			wg.Add(1)
			envId := envelope.EventMeta.Uuid.String()
			_, err := s.client.Create(index, envId, reader)
			if err != nil {
				log.Error().Interface("envelopeId", envId).Err(err).Msg("🔴 could not publish envelope to elasticsearch")
				return err
			} else {
				log.Debug().Interface("🟡 envelopeId", envId).Interface("indexId", index).Msg("published envelope to index")
			}
			defer wg.Done()
		}
	}
	wg.Wait()
	return nil
}

func (s *ElasticsearchSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validIndex, validEnvelopes)
	return err
}

func (s *ElasticsearchSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.invalidIndex, invalidEnvelopes)
	return err
}

func (s *ElasticsearchSink) Close() {
	log.Debug().Msg("🟡 closing elasticsearch sink client")
}
