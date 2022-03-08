package config

type Config struct {
	App         `json:"app"`
	Middleware  `json:"middleware"`
	Snowplow    `json:"snowplow"`
	Generic     `json:"generic"`
	Cloudevents `json:"cloudevents"`
	Sink        `json:"sink"`
	SchemaCache `json:"schemaCache"`
	Tele        `json:"tele"`
}
