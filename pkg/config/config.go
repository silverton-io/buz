package config

type Config struct {
	App       `json:"app"`
	Routing   `json:"routing"`
	Cookie    `json:"cookie"`
	Cors      `json:"cors"`
	Pubsub    `json:"pubsub"`
	Anonymize `json:"anonymize"`
	Cache     `json:"cache"`
	Tele      `json:"tele"`
}
