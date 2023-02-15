// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package mongodb

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Sink struct {
	id                *uuid.UUID
	name              string
	deliveryRequired  bool
	client            *mongo.Client
	validCollection   *mongo.Collection
	invalidCollection *mongo.Collection
}

func (s *Sink) Id() *uuid.UUID {
	return s.id
}

func (s *Sink) Name() string {
	return s.name
}

func (s *Sink) Type() string {
	return "mongodb"
}

func (s *Sink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing mongodb sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	ctx := context.Background()
	opt := options.ClientOptions{
		Hosts: conf.MongoHosts,
	}
	if conf.MongoUser != "" {
		c := options.Credential{
			Username: conf.MongoUser,
			Password: conf.MongoPass,
		}
		opt.Auth = &c
	}
	client, err := mongo.Connect(ctx, &opt)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not connect to mongodb")
	}
	s.client = client
	vCollection := s.client.Database(conf.MongoDbName).Collection(constants.BUZ_VALID_EVENTS)
	iCollection := s.client.Database(conf.MongoDbName).Collection(constants.BUZ_INVALID_EVENTS)
	s.validCollection, s.invalidCollection = vCollection, iCollection
	return nil
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	for _, e := range envelopes {
		payload, err := bson.Marshal(e)
		if err != nil {
			return err
		}
		_, err = s.validCollection.InsertOne(ctx, payload) // FIXME - should batch these and shard them
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sink) Close() {
	log.Debug().Msg("🟡 closing mongodb sink")
}
