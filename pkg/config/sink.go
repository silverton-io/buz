package config

type Sink struct {
	Type                  string   `json:"type"`
	Project               string   `json:"project,omitempty"`
	KafkaBrokers          []string `json:"kakfaBrokers,omitempty"`
	ElasticsearchHosts    []string `json:"elasticsearchHosts,omitempty"`
	ElasticsearchUsername string   `json:"-"`
	ElasticsearchPassword string   `json:"-"`
	ValidEventTopic       string   `json:"validEventTopic,omitempty"`
	InvalidEventTopic     string   `json:"invalidEventTopic,omitempty"`
	ValidUrl              string   `json:"validUrl,omitempty"`
	InvalidUrl            string   `json:"invalidUrl,omitempty"`
	ValidIndex            string   `json:"validIndex,omitempty"`
	InvalidIndex          string   `json:"invalidIndex,omitempty"`
	RelayUrl              string   `json:"relayUrl,omitempty"`
}
