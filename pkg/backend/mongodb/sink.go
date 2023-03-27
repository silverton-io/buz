// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package mongodb

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	client   *mongo.Client
	input    chan []envelope.Envelope
	shutdown chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing mongodb sink")
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	ctx := context.Background()
	opt := options.ClientOptions{
		Hosts: conf.Hosts,
	}
	if conf.User != "" {
		c := options.Credential{
			Username: conf.User,
			Password: conf.Password,
		}
		opt.Auth = &c
	}
	client, err := mongo.Connect(ctx, &opt)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not connect to mongodb")
	}
	s.client = client
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
	return nil
}

func (s *Sink) StartWorker() error {
	err := backendutils.StartSinkWorker(s.input, s.shutdown, s)
	return err
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.input <- envelopes
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope, output string) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	for _, e := range envelopes {
		payload, err := bson.Marshal(e)
		if err != nil {
			return err
		}
		collection := s.client.Database(s.metadata.Name).Collection(output)
		_, err = collection.InsertOne(ctx, payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	s.shutdown <- 1
	return nil
}
