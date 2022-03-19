package config

import "time"

type Sink struct {
	Type                  string        `json:"type"`
	Project               string        `json:"project,omitempty"`
	KafkaBrokers          []string      `json:"kakfaBrokers,omitempty"`
	ElasticsearchHosts    []string      `json:"elasticsearchHosts,omitempty"`
	ElasticsearchUsername string        `json:"elasticsearchUsername,omitempty"`
	ElasticsearchPassword string        `json:"elasticsearchPassword,omitempty"`
	ValidEventTopic       string        `json:"validEventTopic,omitempty"`
	InvalidEventTopic     string        `json:"invalidEventTopic,omitempty"`
	ValidUrl              string        `json:"validUrl,omitempty"`
	InvalidUrl            string        `json:"invalidUrl,omitempty"`
	ValidIndex            string        `json:"validIndex,omitempty"`
	InvalidIndex          string        `json:"invalidIndex,omitempty"`
	RelayUrl              string        `json:"relayUrl,omitempty"`
	BufferByteThreshold   int           `json:"bufferByteThreshold,omitempty"`
	BufferRecordThreshold int           `json:"bufferRecordThreshold,omitempty"`
	BufferDelayThreshold  time.Duration `json:"bufferDelayThreshold,omitempty"`
}
