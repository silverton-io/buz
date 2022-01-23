package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

var ValidClient *pubsub.Client
var ValidTopic *pubsub.Topic
var InvalidClient *pubsub.Client
var InvalidTopic *pubsub.Topic

func buildPubsubClients() {
	ctx := context.Background()
	validClient, err := pubsub.NewClient(ctx, Config.pubsubProjectName)
	invalidClient, err := pubsub.NewClient(ctx, Config.pubsubProjectName)
	if err != nil {
		log.Fatalf("Failure creating a pubsub client %v", err)
	}
	validTopic := validClient.Topic(Config.pubsubValidTopicName)
	invalidTopic := invalidClient.Topic(Config.pubsubInvalidTopicName)

	ValidClient = validClient
	ValidTopic = validTopic

	InvalidClient = invalidClient
	InvalidTopic = invalidTopic

	defer InvalidClient.Close()
	defer ValidClient.Close()
}
