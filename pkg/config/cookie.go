package config

type Cookie struct {
	Name    string `json:"name"`
	Secure  bool   `json:"secure"`
	TtlDays int    `json:"ttlDays"`
	Domain  string `json:"domain"`
	Path    string `json:"path"`
}
