package main

type AppConfig struct {
	environment           string
	port                  string
	includeStandardRoutes bool
	// Custom routing
	postPath     string
	getPath      string
	redirectPath string
	// Internal functionality
	cookieName     string
	cookieDomain   string
	cookiePath     string
	cookieTtlDays  int
	extractCookies string
	// CORS
	corsAllowedOrigins []string
	corsMaxAge         int
	// Operational
	// schemaCacheMemoryMax    int
	// schemaCacheTtlSeconds   int
	// schemaCachePurgeSeconds int
	// schemaCachePurgeToken   string
}

var Config = AppConfig{
	environment:           "dev",
	port:                  "8080",
	includeStandardRoutes: true,
	postPath:              "com.snowplowanalytics.snowplow/tp2",
	getPath:               "/i",
	redirectPath:          "r/tp2",
	cookieName:            "sp-nuid",
	cookieDomain:          "localhost",
	cookiePath:            "/",
	cookieTtlDays:         365,
	corsAllowedOrigins:    []string{"*"},
	corsMaxAge:            86400,
}
