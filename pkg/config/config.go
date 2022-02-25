package config

type Config struct {
	App         `json:"app"`
	Snowplow    `json:"snowplow"`
	Generic     `json:"generic"`
	Cloudevents `json:"cloudevents"`
	Cookie      `json:"cookie"`
	Cors        `json:"cors"`
	Sink        `json:"sink"`
	SchemaCache `json:"schemaCache"`
	Stats       `json:"stats"`
	Tele        `json:"tele"`
}
