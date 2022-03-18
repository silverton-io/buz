package config

type Config struct {
	App         `json:"app"`
	Middleware  `json:"middleware"`
	Inputs      `json:"inputs"`
	Sink        `json:"sink"`
	SchemaCache `json:"schemaCache"`
	Squawkbox   `json:"squawkBox"`
	Tele        `json:"tele"`
}
