package config

type Tele struct {
	Enabled     bool   `json:"enabled,omitempty"`
	HeartbeatMs int    `json:"heartbeatMs"`
	Host        string `json:"host,omitempty"`
}
