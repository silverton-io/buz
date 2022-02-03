package config

type Config struct {
	App             `mapstructure:"app"`
	AdvancingCookie `mapstructure:"advancingCookie"`
	Cors            `mapstructure:"cors"`
	Pubsub          `mapstructure:"pubsub"`
}
