package config

type Sink struct {
	Name         string   `json:"name"`
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
	// Database
	DbHost       string `json:"-"`
	DbPort       uint16 `json:"-"`
	DbName       string `json:"-"`
	DbUser       string `json:"-"`
	DbPass       string `json:"-"`
	ValidTable   string `json:"validTable,omitempty"`
	InvalidTable string `json:"invalidTable,omitempty"`
	// Pubnub
	ValidChannel   string `json:"validChannel,omitempty"`
	InvalidChannel string `json:"invalidChannel,omitempty"`
	PubnubPubKey   string `json:"pubnubPubKey,omitempty"`
	PubnubSubKey   string `json:"pubnubSubKey,omitempty"`
	// Mongodb
	MongoHosts        []string `json:"mongoHosts,omitempty"`
	MongoDbPort       string   `json:"mongoDbPort,omitempty"`
	MongoDbName       string   `json:"mongoDbName,omitempty"`
	MongoDbUser       string   `json:"-"`
	MongoDbPass       string   `json:"-"`
	ValidCollection   string   `json:"validCollection,omitempty"`
	InvalidCollection string   `json:"invalidCollection,omitempty"`
}
