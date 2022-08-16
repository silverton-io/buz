// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package cache

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSchemaDocument struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Contents string             `bson:"contents"`
}

type MongodbSchemaCacheBackend struct {
	client             *mongo.Client
	registryCollection *mongo.Collection
}

func (b *MongodbSchemaCacheBackend) Initialize(conf config.Backend) error {
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
	b.client = client
	registryCollection := b.client.Database(conf.MongoDbName).Collection(conf.RegistryCollection)
	b.registryCollection = registryCollection
	return nil
}

func (b *MongodbSchemaCacheBackend) GetRemote(schema string) (contents []byte, err error) {
	ctx := context.Background()
	var doc = MongoSchemaDocument{}
	err = b.registryCollection.FindOne(ctx, bson.M{"name": schema}).Decode(&doc)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not decode document")
		return nil, err
	}
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not marshal document")
		return nil, err
	}
	return []byte(doc.Contents), nil
}

func (b *MongodbSchemaCacheBackend) Close() {
	log.Info().Msg("🟢 closing mongodb schema cache backend")
}
