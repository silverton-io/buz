package config

type Squawkbox struct {
	Enabled         bool   `json:"enabled"`
	CloudeventsPath string `json:"cloudeventsPath"`
	SnowplowPath    string `json:"snowplowPath"`
	GenericPath     string `json:"genericPath"`
}
