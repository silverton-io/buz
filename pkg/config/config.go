package config

type Config struct {
	App         `json:"app"`
	Middleware  `json:"middleware"`
	Inputs      `json:"inputs"`
	SchemaCache `json:"schemaCache"`
	Manifold    `json:"manifold"`
	Sinks       []Sink `json:"sinks"`
	Squawkbox   `json:"squawkBox"`
	Privacy     `json:"privacy"`
	Tele        `json:"tele"`
}
