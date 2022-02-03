package config

type Cors struct {
	AllowOrigin      []string `mapstructure:"allowOrigin"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	AllowMethods     []string `mapstructure:"allowMethods"`
	MaxAge           int      `mapstructure:"maxAge"`
}
