package config

import "time"

type Pubsub struct {
	Project               string        `json:"project"`
	ValidEventTopic       string        `json:"validEventTopic"`
	InvalidEventTopic     string        `json:"invalidEventTopic"`
	BufferByteThreshold   int           `json:"bufferByteThreshold"`
	BufferRecordThreshold int           `json:"bufferRecordThreshold"`
	BufferDelayThreshold  time.Duration `json:"bufferDelayThreshold"`
}
