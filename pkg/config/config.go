package config

type Config struct {
	App         `json:"app"`
	Middleware  `json:"middleware"`
	Snowplow    `json:"snowplow"`
	Generic     `json:"generic"`
	Cloudevents `json:"cloudevents"`
	Webhook     `json:"webhook"`
	Sink        `json:"sink"`
	SchemaCache `json:"schemaCache"`
	Squawkbox   `json:"squawkBox"`
	Relay       `json:"relay"`
	Tele        `json:"tele"`
}
