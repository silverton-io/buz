package config

type Config struct {
	App         `json:"app"`
	Routing     `json:"routing"`
	Cookie      `json:"cookie"`
	Cors        `json:"cors"`
	Forwarder   `json:"forwarder"`
	Anonymize   `json:"anonymize"`
	SchemaCache `json:"schemaCache"`
	Tele        `json:"tele"`
}