package main

type AppConfig struct {
	environment                   string
	port                          string
	includeStandardSnowplowRoutes bool
	snowplowPostPath              string
	snowplowGetPath               string
	snowplowRedirectPath          string
	validateSnowplowEvents        bool
}

type PubsubConfig struct {
	projectName       string
	validEventTopic   string
	invalidEventTopic string
	// bufferByteThreshold           int
	// bufferCountThreshold          int
	// bufferDelayThreshold          time.Duration
}
