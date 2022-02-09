package config

type Backend struct {
	Type     string `json:"type"`
	Location string `json:"location"`
	Path     string `json:"path"`
}

type Cache struct {
	Backend      `json:"backend"`
	TtlSeconds   int `json:"ttlSeconds"`
	MaxSizeBytes int `json:"maxSizeBytes"`
}
