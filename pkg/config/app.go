package config

type App struct {
	Env      string `json:"env"`
	Port     string `json:"port"`
	Mode     string `json:"mode"`
	LogLevel string `json:"logLevel"`
}
