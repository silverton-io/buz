package config

type AdvancingCookie struct {
	CookieName      string `mapstructure:"cookieName"`
	UseSecureCookie bool   `mapstructure:"useSecureCookie"`
	CookieTtlDays   int    `mapstructure:"cookieTtlDays"`
	CookiePath      string `mapstructure:"cookiePath"`
	CookieDomain    string `mapstructure:"cookieDomain"`
}
