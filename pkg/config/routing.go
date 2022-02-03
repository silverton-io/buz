package config

type Routing struct {
	IncludeStandardRoutes bool   `json:"includeStandardRoutes"`
	GetPath               string `json:"getPath"`
	PostPath              string `json:"postPath"`
	RedirectPath          string `json:"redirectPath"`
}
