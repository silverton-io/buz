package config

type Config struct {
	App         `json:"app"`
	Middleware  `json:"middleware"`
	Inputs      `json:"inputs"`
	SchemaCache `json:"schemaCache"`
	Manifold    `json:"manifold"`
	Sink        `json:"sink"`
	Squawkbox   `json:"squawkBox"`
	Tele        `json:"tele"`
}
