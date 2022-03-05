package config

type RateLimiter struct {
	Period string `json:"period"`
	Limit  int64  `json:"limit"`
}
