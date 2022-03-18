package config

import "time"

type Sink struct {
	Type                  string        `json:"type"`
	Project               string        `json:"project,omitempty"`
	Brokers               []string      `json:"brokers,omitempty"`
	ValidEventTopic       string        `json:"validEventTopic"`
	InvalidEventTopic     string        `json:"invalidEventTopic"`
	ValidUrl              string        `json:"validUrl"`
	InvalidUrl            string        `json:"invalidUrl"`
	RelayUrl              string        `json:"relayUrl"`
	BufferByteThreshold   int           `json:"bufferByteThreshold,omitempty"`
	BufferRecordThreshold int           `json:"bufferRecordThreshold,omitempty"`
	BufferDelayThreshold  time.Duration `json:"bufferDelayThreshold,omitempty"`
}
