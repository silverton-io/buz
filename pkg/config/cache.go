package config

type SchemaCacheBackend struct {
	Type   string `json:"type"`
	Region string `json:"region"`
	Bucket string `json:"bucket"`
	Path   string `json:"path"`
}

type SchemaCache struct {
	SchemaCacheBackend `json:"schemaCacheBackend"`
	TtlSeconds         int `json:"ttlSeconds"`
	MaxSizeBytes       int `json:"maxSizeBytes"`
}
