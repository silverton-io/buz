package config

type SelfDescribingConfig struct {
	RootKey   string `json:"rootKey"`
	SchemaKey string `json:"schemaKey"`
	DataKey   string `json:"dataKey"`
}

type Generic struct {
	Enabled  bool                 `json:"enabled"`
	Path     string               `json:"path"`
	Contexts SelfDescribingConfig `json:"contexts"`
	Payload  SelfDescribingConfig `json:"payload"`
}
