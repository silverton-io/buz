package config

type App struct {
	Version       string `json:"version"`
	Env           string `json:"env"`
	Mode          string `json:"mode"`
	Port          string `json:"port"`
	TrackerDomain string `json:"trackerDomain"`
	Stats         `json:"stats"`
	Timeout       `json:"timeout"`
}
