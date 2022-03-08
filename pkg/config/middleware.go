package config

type Middleware struct {
	Timeout       `json:"timeout"`
	RateLimiter   `json:"rateLimiter"`
	Cookie        `json:"cookie"`
	Cors          `json:"cors"`
	RequestLogger `json:"requestLogger"`
	Yeet          `json:"yeet"`
}
