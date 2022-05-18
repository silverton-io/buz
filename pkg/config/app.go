package config

type App struct {
	Version       string `json:"version"`
	Name          string `json:"name"`
	Env           string `json:"env"`
	Mode          string `json:"mode"`
	Port          string `json:"port"`
	TrackerDomain string `json:"trackerDomain"`
	Health        `json:"health"`
	Stats         `json:"stats"`
}

type Stats struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type Health struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}
