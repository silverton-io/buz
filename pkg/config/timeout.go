package config

type Timeout struct {
	Enabled bool `json:"enabled"`
	Ms      int  `json:"ms"`
}
