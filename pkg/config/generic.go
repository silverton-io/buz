package config

type Generic struct {
	Enabled       bool   `json:"enabled"`
	PostPath      string `json:"postPath"`
	BatchPostPath string `json:"batchPostPath"`
	ContextsKey   string `json:"contextsKey"`
	PayloadKey    string `json:"payloadKey"`
}
