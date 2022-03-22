package config

type Sinks struct {
	Sinks []Sink `json:"sinks"`
}

type Sink struct {
	Type         string   `json:"type"`
	Project      string   `json:"project,omitempty"`
	KafkaBrokers []string `json:"kakfaBrokers,omitempty"`
	// Kafka, Pubsub
	ValidEventTopic   string `json:"validEventTopic,omitempty"`
	InvalidEventTopic string `json:"invalidEventTopic,omitempty"`
	// Relay, HTTP/S, etc
	ValidUrl   string `json:"validUrl,omitempty"`
	InvalidUrl string `json:"invalidUrl,omitempty"`
	// Elasticsearch
	ValidIndex            string   `json:"validIndex,omitempty"`
	InvalidIndex          string   `json:"invalidIndex,omitempty"`
	ElasticsearchHosts    []string `json:"elasticsearchHosts,omitempty"`
	ElasticsearchUsername string   `json:"-"`
	ElasticsearchPassword string   `json:"-"`
	// Honeypot relay
	RelayUrl string `json:"relayUrl,omitempty"`
	// File
	ValidFile   string `json:"validFile,omitempty"`
	InvalidFile string `json:"invalidFile,omitempty"`
}
