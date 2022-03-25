package config

type Cloudevents struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}
