package config

type Snowplow struct {
	Enabled               bool   `json:"enabled"`
	StandardRoutesEnabled bool   `json:"standardRoutesEnabled"`
	OpenRedirectsEnabled  bool   `json:"openRedirectsEnabled"`
	GetPath               string `json:"getPath"`
	PostPath              string `json:"postPath"`
	RedirectPath          string `json:"redirectPath"`
}
