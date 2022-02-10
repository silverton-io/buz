package config

import "time"

type Forwarder struct {
	Type                  string        `json:"type"`
	Project               string        `json:"project"`
	ValidEventTopic       string        `json:"validEventTopic"`
	InvalidEventTopic     string        `json:"invalidEventTopic"`
	BufferByteThreshold   int           `json:"bufferByteThreshold"`
	BufferRecordThreshold int           `json:"bufferRecordThreshold"`
	BufferDelayThreshold  time.Duration `json:"bufferDelayThreshold"`
}
