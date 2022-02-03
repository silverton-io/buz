package config

type App struct {
	Env                   string `mapstructure:"env"`
	Port                  string `mapstructure:"port"`
	IncludeStandardRoutes bool   `mapstructure:"includeStandardRoutes"`
}
