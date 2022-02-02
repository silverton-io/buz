package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

var ValidClient *pubsub.Client
var InvalidClient *pubsub.Client

func buildPubsubClients() {
	ctx := context.Background()
	validClient, err := pubsub.NewClient(ctx, Config.pubsubProjectName)
	// invalidClient, err := pubsub.NewClient(ctx, Config.pubsubProjectName)
	if err != nil {
		log.Fatalf("Failure creating a pubsub client %v", err)
	}

	ValidClient = validClient
	// InvalidClient = invalidClient

	defer ValidClient.Close()
	// defer InvalidClient.Close()
}

func publishEvents(topicName string, events []Event) {
	ctx := context.Background()
	topic := ValidClient.Topic(topicName)

	for _, event := range events {
		fmt.Println("Made it here")
		marshaledEvent, _ := json.Marshal(event)
		result := topic.Publish(ctx, &pubsub.Message{Data: marshaledEvent})
		fmt.Println("And then here")
		id, err := result.Get(ctx)
		if err != nil {
			fmt.Printf("Something bad happened %s", err)
		}
		fmt.Println(id)
		// wg.Add(1)
		// go func(res *pubsub.PublishResult) {
		// 	defer wg.Done()
		// 	id, err := res.Get(ctx)
		// 	if err != nil {
		// 		fmt.Printf("Error when publishing event: %s", err)
		// 	}
		// 	fmt.Printf("Message published. Id: %s", id)
		// }(result)
	}
	fmt.Println("And then here")
	// wg.Wait()
	fmt.Println("And finally returned")
}

func publishValidEvents(events []Event) {
	publishEvents(Config.pubsubValidTopicName, events)
}

// func publishInvalidEvents(events []Event) {
// 	publishEvents(InvalidTopic, events)
// }
