package config

type SelfDescribingConfig struct {
	RootKey   string `json:"rootKey"`
	SchemaKey string `json:"schemaKey"`
	DataKey   string `json:"dataKey"`
}

type Generic struct {
	Enabled       bool                 `json:"enabled"`
	PostPath      string               `json:"postPath"`
	BatchPostPath string               `json:"batchPostPath"`
	Contexts      SelfDescribingConfig `json:"contexts"`
	Event         SelfDescribingConfig `json:"event"`
}
