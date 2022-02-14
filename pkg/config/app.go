package config

type App struct {
	Version string `json:"version"`
	Env     string `json:"env"`
	Port    string `json:"port"`
	Mode    string `json:"mode"`
}
