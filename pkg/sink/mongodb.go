package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbSink struct {
	id                *uuid.UUID
	name              string
	client            *mongo.Client
	validCollection   *mongo.Collection
	invalidCollection *mongo.Collection
}

func (s *MongodbSink) Id() *uuid.UUID {
	return s.id
}

func (s *MongodbSink) Name() string {
	return s.name
}

func (s *MongodbSink) Initialize(conf config.Sink) error {
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	ctx := context.Background()
	opt := options.ClientOptions{
		Hosts: conf.MongoHosts,
	}
	if conf.MongoDbUser != "" {
		c := options.Credential{
			Username: conf.MongoDbUser,
			Password: conf.MongoDbPass,
		}
		opt.Auth = &c
	}
	client, err := mongo.Connect(ctx, &opt)
	if err != nil {
		log.Error().Err(err).Msg("could not connect to mongo")
	}
	s.client = client
	vCollection := s.client.Database(conf.MongoDbName).Collection(conf.ValidCollection)
	iCollection := s.client.Database(conf.MongoDbName).Collection(conf.InvalidCollection)
	s.validCollection, s.invalidCollection = vCollection, iCollection
	return nil
}

func (s *MongodbSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	for _, e := range envelopes {
		payload, err := bson.Marshal(e)
		if err != nil {
			return err
		}
		_, err = s.validCollection.InsertOne(ctx, payload) // FIXME - should batch these
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MongodbSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	for _, e := range envelopes {
		payload, err := bson.Marshal(e)
		if err != nil {
			return err
		}
		_, err = s.invalidCollection.InsertOne(ctx, payload) // FIXME - should batch these
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MongodbSink) Close() {
	log.Debug().Msg("closing mongodb sink")
}
