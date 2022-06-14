package config

type Middleware struct {
	Timeout       `json:"timeout"`
	RateLimiter   `json:"rateLimiter"`
	Cookie        `json:"cookie"`
	Cors          `json:"cors"`
	RequestLogger `json:"requestLogger"`
	Yeet          `json:"yeet"`
}

type Timeout struct {
	Enabled bool `json:"enabled"`
	Ms      int  `json:"ms"`
}

type RateLimiter struct {
	Enabled bool   `json:"enabled"`
	Period  string `json:"period"`
	Limit   int64  `json:"limit"`
}

type Cookie struct {
	Enabled  bool   `json:"enabled"`
	Name     string `json:"name"`
	Secure   bool   `json:"secure"`
	TtlDays  int    `json:"ttlDays"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	SameSite string `json:"sameSite"`
	Fallback string `json:"fallback"`
}

type Cors struct {
	Enabled          bool     `json:"enabled"`
	AllowOrigin      []string `json:"allowOrigin"`
	AllowCredentials bool     `json:"allowCredentials"`
	AllowMethods     []string `json:"allowMethods"`
	MaxAge           int      `json:"maxAge"`
}

type RequestLogger struct {
	Enabled bool `json:"enabled"`
}

type Yeet struct {
	Enabled bool `json:"enabled"`
}
