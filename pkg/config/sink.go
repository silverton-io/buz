// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package config

type Sink struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	DeliveryRequired bool   `json:"deliveryRequired"`
	Fanout           bool   `json:"fanout"`
	Deadletter       string `json:"deadletter"`
	// GCP
	Project string `json:"project,omitempty"`
	// Kafka
	Brokers []string `json:"kakfaBrokers,omitempty"`
	// HTTP and API
	Url    string `json:"url,omitempty"`
	ApiKey string `json:"-"`
	Region string `json:"-"`
	// Database
	Hosts    []string `json:"-"`
	Port     uint16   `json:"-"`
	Database string   `json:"-"`
	User     string   `json:"-"`
	Password string   `json:"-"`
	// Pubnub
	PubnubPubKey string `json:"pubnubPubKey,omitempty"`
	PubnubSubKey string `json:"pubnubSubKey,omitempty"`
}
