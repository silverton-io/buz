package config

type SchemaCacheBackend struct {
	Type     string `json:"type"`
	Location string `json:"location"`
	Path     string `json:"path"`
}

type SchemaCache struct {
	SchemaCacheBackend `json:"schemaCacheBackend"`
	TtlSeconds         int `json:"ttlSeconds"`
	MaxSizeBytes       int `json:"maxSizeBytes"`
}
