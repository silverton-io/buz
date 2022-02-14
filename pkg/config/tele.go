package config

type Tele struct {
	Enable bool   `json:"enable,omitempty"`
	Host   string `json:"host,omitempty"`
}
