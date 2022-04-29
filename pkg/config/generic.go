package config

type SingleLayerConfig struct {
	RootKey string `json:"rootKey"`
}

type DoubleLayerConfig struct {
	RootKey   string `json:"rootKey"`
	SchemaKey string `json:"schemaKey"`
	DataKey   string `json:"dataKey"`
}

type Generic struct {
	Enabled  bool              `json:"enabled"`
	Path     string            `json:"path"`
	Contexts SingleLayerConfig `json:"contexts"`
	Payload  DoubleLayerConfig `json:"payload"`
}
