package config

type Anonymize struct {
	Ip     bool `json:"ip"`
	UserId bool `json:"userId"`
}

type Snowplow struct {
	Enabled               bool   `json:"enabled"`
	StandardRoutesEnabled bool   `json:"standardRoutesEnabled"`
	OpenRedirectsEnabled  bool   `json:"openRedirectsEnabled"`
	GetPath               string `json:"getPath"`
	PostPath              string `json:"postPath"`
	RedirectPath          string `json:"redirectPath"`
	Anonymize             `json:"anonymize"`
}
