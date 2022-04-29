package config

type GenericRootConfig struct {
	RootKey string `json:"rootKey"`
}

type GenericRootAndChildConfig struct {
	RootKey   string `json:"rootKey"`
	SchemaKey string `json:"schemaKey"`
	DataKey   string `json:"dataKey"`
}

type Generic struct {
	Enabled  bool                      `json:"enabled"`
	Path     string                    `json:"path"`
	Contexts GenericRootConfig         `json:"contexts"`
	Payload  GenericRootAndChildConfig `json:"payload"`
}
