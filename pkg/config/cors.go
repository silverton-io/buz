package config

type Cors struct {
	Enabled          bool     `json:"enabled"`
	AllowOrigin      []string `json:"allowOrigin"`
	AllowCredentials bool     `json:"allowCredentials"`
	AllowMethods     []string `json:"allowMethods"`
	MaxAge           int      `json:"maxAge"`
}
