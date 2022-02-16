package config

type Config struct {
	App         `json:"app"`
	Snowplow    `json:"snowplow"`
	Generic     `json:"generic"`
	Cookie      `json:"cookie"`
	Cors        `json:"cors"`
	Forwarder   `json:"forwarder"`
	Anonymize   `json:"anonymize"`
	SchemaCache `json:"schemaCache"`
	Stats       `json:"stats"`
	Tele        `json:"tele"`
}
