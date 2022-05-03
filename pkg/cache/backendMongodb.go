package cache

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
		log.Error().Err(err).Msg("could not connect to mongodb")
	}
	b.client = client
	registryCollection := b.client.Database(conf.MongoDbName).Collection(conf.RegistryCollection)
	b.registryCollection = registryCollection
	return nil
}

func (b *MongodbSchemaCacheBackend) GetRemote(conf config.Backend) (contents []byte, err error) {}

func (b *MongodbSchemaCacheBackend) Close() {
	log.Info().Msg("closing mongodb schema cache backend")
}
