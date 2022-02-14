package config

type App struct {
	Version    string `json:"version"`
	InstanceId string `json:"instanceId"`
	Env        string `json:"env"`
	Port       string `json:"port"`
	Mode       string `json:"mode"`
}
