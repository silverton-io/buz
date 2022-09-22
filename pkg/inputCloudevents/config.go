package inputcloudevents

type CloudeventsConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}
