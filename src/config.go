package main

type AppConfig struct {
	environment string
	port        string
	// Routing
	includeStandardSnowplowRoutes bool
	snowplowPostPath              string
	snowplowGetPath               string
	snowplowRedirectPath          string
	openPostPath                  string
	// Identification
	cookieName     string
	cookieDomain   string
	cookiePath     string
	cookieTtlDays  int
	extractCookies string
	// CORS
	corsAllowedOrigins []string
	corsMaxAge         int
	// Destination
	pubsubProjectName      string
	pubsubValidTopicName   string
	pubsubInvalidTopicName string
	// Event Validation
	validateSnowplowEvents bool
	// Operational
	// schemaCacheMemoryMax    int
	// schemaCacheTtlSeconds   int
	// schemaCachePurgeSeconds int
	// schemaCachePurgeToken   string
}

var Config = AppConfig{
	environment: "dev",
	port:        "8080",
	// Routing
	includeStandardSnowplowRoutes: true,
	snowplowPostPath:              "com.snowplowanalytics.snowplow/tp2",
	snowplowGetPath:               "/i",
	snowplowRedirectPath:          "r/tp2",
	// Identification
	cookieName:    "sp-nuid",
	cookieDomain:  "localhost",
	cookiePath:    "/",
	cookieTtlDays: 365,
	// CORS
	corsAllowedOrigins: []string{"*"},
	corsMaxAge:         86400,
	// Destination
	pubsubProjectName:      "neat-dispatch-338321",
	pubsubValidTopicName:   "test-topic",
	pubsubInvalidTopicName: "test-topic",
	// Validation
	validateSnowplowEvents: false,
}
