package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
)

func publishEvent(ctx context.Context, topic *pubsub.Topic, event Event) {
	payload, _ := json.Marshal(event)
	msg := &pubsub.Message{
		Data: payload,
	}
	result := topic.Publish(ctx, msg)
	id, err := result.Get(ctx)
	if err != nil {
		fmt.Printf("FIXME! something bad happened %v\n", err)
	}
	fmt.Printf("message %v sent to pubsub\n", id)
}

func publishEvents(ctx context.Context, topic *pubsub.Topic, events []Event) {
	var wg sync.WaitGroup
	for _, event := range events {
		payload, _ := json.Marshal(event)
		msg := &pubsub.Message{
			Data: payload,
		}
		result := topic.Publish(ctx, msg)
		wg.Add(1)
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			_, err := res.Get(ctx)
			if err != nil {
				fmt.Printf("FIXME! something bad happened %v\n", err)
			}
		}(result)
	}
	wg.Wait()
}
