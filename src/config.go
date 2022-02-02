package main

import "time"

type AppConfig struct {
	environment                   string
	port                          string
	includeStandardSnowplowRoutes bool
	snowplowPostPath              string
	snowplowGetPath               string
	snowplowRedirectPath          string
	cookieName                    string
	cookieDomain                  string
	cookiePath                    string
	cookieTtlDays                 int
	corsAllowedOrigins            []string
	corsMaxAge                    int
	pubsubProjectName             string
	pubsubValidTopicName          string
	pubsubInvalidTopicName        string
	bufferByteThreshold           int
	bufferCountThreshold          int
	bufferDelayThreshold          time.Duration
	validateSnowplowEvents        bool
	// Operational
	// schemaCacheMemoryMax    int
	// schemaCacheTtlSeconds   int
	// schemaCachePurgeSeconds int
	// schemaCachePurgeToken   string
}

var Config = AppConfig{
	environment:                   "dev",
	port:                          "8080",
	includeStandardSnowplowRoutes: true,
	snowplowPostPath:              "com.snowplowanalytics.snowplow/tp2",
	snowplowGetPath:               "/i",
	snowplowRedirectPath:          "r/tp2",
	cookieName:                    "sp-nuid",
	cookieDomain:                  "localhost",
	cookiePath:                    "/",
	cookieTtlDays:                 365,
	corsAllowedOrigins:            []string{"*"},
	corsMaxAge:                    86400,
	pubsubProjectName:             "neat-dispatch-338321",
	pubsubValidTopicName:          "test-valid-topic",
	pubsubInvalidTopicName:        "test-invalid-topic",
	bufferByteThreshold:           0,
	bufferCountThreshold:          0,
	bufferDelayThreshold:          1 * time.Millisecond,
	validateSnowplowEvents:        false,
}
