package config

type Middleware struct {
	Timeout       `json:"timeout"`
	RateLimiter   `json:"rateLimiter"`
	Cookie        `json:"cookie"`
	Cors          `json:"cors"`
	RequestLogger `json:"jsonLogger"`
	Yeet          `json:"yeet"`
}
