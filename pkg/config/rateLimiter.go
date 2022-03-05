package config

type RateLimiter struct {
	Enabled bool   `json:"enabled"`
	Period  string `json:"period"`
	Limit   int64  `json:"limit"`
}
