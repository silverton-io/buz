package config

import "time"

type Forwarder struct {
	Type                  string        `json:"type"`
	Project               string        `json:"project,omitempty"`
	Brokers               []string      `json:"brokers,omitempty"`
	ValidEventTopic       string        `json:"validEventTopic"`
	InvalidEventTopic     string        `json:"invalidEventTopic"`
	BufferByteThreshold   int           `json:"bufferByteThreshold,omitempty"`
	BufferRecordThreshold int           `json:"bufferRecordThreshold,omitempty"`
	BufferDelayThreshold  time.Duration `json:"bufferDelayThreshold,omitempty"`
}
