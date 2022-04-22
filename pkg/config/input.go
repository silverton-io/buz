package config

type Inputs struct {
	Snowplow    `json:"snowplow"`
	Cloudevents `json:"cloudevents"`
	Generic     `json:"generic"`
	Webhook     `json:"webhook"`
	Pixel       `json:"pixel"`
	Relay       `json:"relay"`
}
