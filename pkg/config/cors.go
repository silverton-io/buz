package config

type Cors struct {
	AllowOrigin      []string `json:"allowOrigin"`
	AllowCredentials bool     `json:"allowCredentials"`
	AllowMethods     []string `json:"allowMethods"`
	MaxAge           int      `json:"maxAge"`
}
