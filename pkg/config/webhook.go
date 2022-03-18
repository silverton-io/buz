package config

type Webhook struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}
