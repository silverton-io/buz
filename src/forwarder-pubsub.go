package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

var PubsubClient *pubsub.Client
var PubsubTopic *pubsub.Topic

func buildPubsubClient() {
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, Config.pubsubProjectName)
	if err != nil {
		log.Fatalf("Failure creating a pubsub topic %v", err)
	}
	pubsubTopic := pubsubClient.Topic(Config.pubsubTopicName)
	pubsubTopic.PublishSettings.NumGoroutines = 1

	PubsubClient = pubsubClient
	PubsubTopic = pubsubTopic
}
