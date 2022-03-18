package config

type Inputs struct {
	Snowplow    `json:"snowplow"`
	Cloudevents `json:"cloudevents"`
	Generic     `json:"generic"`
	Webhook     `json:"webhook"`
	Relay       `json:"relay"`
}
