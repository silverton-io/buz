package config

type Routing struct {
	DisableStandardRoutes bool   `json:"disableStandardRoutes"`
	DisableOpenRedirect   bool   `json:"disableOpenRedirect"`
	GetPath               string `json:"getPath"`
	PostPath              string `json:"postPath"`
	RedirectPath          string `json:"redirectPath"`
}
