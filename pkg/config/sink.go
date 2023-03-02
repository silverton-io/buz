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
	// Database
	DbHosts []string `json:"-"`
	DbPort  uint16   `json:"-"`
	DbName  string   `json:"-"`
	DbUser  string   `json:"-"`
	DbPass  string   `json:"-"`
	// Pubnub
	PubnubPubKey string `json:"pubnubPubKey,omitempty"`
	PubnubSubKey string `json:"pubnubSubKey,omitempty"`
	// Indicative
	IndicativeApiKey string `json:"-"`
	// Amplitude
	AmplitudeApiKey string `json:"-"`
	AmplitudeRegion string `json:"amplitudeRegion,omitempty"`
}
