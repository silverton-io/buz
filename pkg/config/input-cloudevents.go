package config

type Cloudevents struct {
	Enabled       bool   `json:"enabled"`
	PostPath      string `json:"postPath"`
	BatchPostPath string `json:"batchPostPath"`
}
