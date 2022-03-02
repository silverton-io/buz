package config

type Purge struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type SchemaCacheBackend struct {
	Type   string `json:"type"`
	Region string `json:"region,omitempty"`
	Bucket string `json:"bucket,omitempty"`
	Host   string `json:"host,omitempty"`
	Path   string `json:"path"`
}

type SchemaCache struct {
	SchemaCacheBackend `json:"schemaCacheBackend"`
	TtlSeconds         int `json:"ttlSeconds"`
	MaxSizeBytes       int `json:"maxSizeBytes"`
	Purge              `json:"purge"`
}
