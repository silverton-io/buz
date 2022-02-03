package config

type Pubsub struct {
	ProjectName       string `mapstructure:"projectName"`
	ValidEventTopic   string `mapstructure:"validEventTopic"`
	InvalidEventTopic string `mapstructure:"invalidEventTopic"`
	// bufferByteThreshold           int
	// bufferCountThreshold          int
	// bufferDelayThreshold          time.Duration
}
