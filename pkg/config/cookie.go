package config

type Cookie struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
	Secure  bool   `json:"secure"`
	TtlDays int    `json:"ttlDays"`
	Domain  string `json:"domain"`
	Path    string `json:"path"`
}
