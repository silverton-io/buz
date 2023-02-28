// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package config

type Sink struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	DeliveryRequired bool   `json:"deliveryRequired"`
	Fanout           bool   `json:"fanout"`
	// Pub/Sub
	Project string `json:"project,omitempty"`
	// Kafka
	KafkaBrokers []string `json:"kakfaBrokers,omitempty"`
	// NATS
	NatsHost string `json:"-"`
	NatsUser string `json:"-"`
	NatsPass string `json:"-"`
	// HTTP
	Url string `json:"url,omitempty"`
	// Elasticsearch
	ElasticsearchHosts    []string `json:"elasticsearchHosts,omitempty"`
	ElasticsearchUsername string   `json:"-"`
	ElasticsearchPassword string   `json:"-"`
	// Postgres Database
	DbHost string `json:"-"`
	DbPort uint16 `json:"-"`
	DbName string `json:"-"`
	DbUser string `json:"-"`
	DbPass string `json:"-"`
	// Pubnub
	PubnubPubKey string `json:"pubnubPubKey,omitempty"`
	PubnubSubKey string `json:"pubnubSubKey,omitempty"`
	// Mongodb
	MongoHosts  []string `json:"mongoHosts,omitempty"`
	MongoPort   string   `json:"mongoDbPort,omitempty"`
	MongoDbName string   `json:"mongoDbName,omitempty"`
	MongoUser   string   `json:"-"`
	MongoPass   string   `json:"-"`
	// Indicative
	IndicativeApiKey string `json:"-"`
	// Amplitude
	AmplitudeApiKey string `json:"-"`
	AmplitudeRegion string `json:"amplitudeRegion,omitempty"`
}
