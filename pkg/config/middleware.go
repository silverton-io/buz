// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package config

type Middleware struct {
	Timeout       `json:"timeout"`
	RateLimiter   `json:"rateLimiter"`
	Identity      `json:"identity"`
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

type Identity struct {
	Cookie   IdentityCookie `json:"cookie"`
	Fallback string         `json:"fallback"`
}

type IdentityCookie struct {
	Enabled  bool   `json:"enabled"`
	Name     string `json:"name"`
	Secure   bool   `json:"secure"`
	TtlDays  int    `json:"ttlDays"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	SameSite string `json:"sameSite"`
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
