package config

type Privacy struct {
	Anonymize  `json:"anonymize"`
	DoNotTrack `json:"doNotTrack"`
}

type DoNotTrack struct {
	Enabled bool `json:"enabled"`
}

type Anonymize struct {
	Device `json:"device"`
	User   `json:"user"`
}

type Device struct {
	Ip        bool `json:"ip"`
	Useragent bool `json:"useragent"`
}

type User struct {
	Id bool `json:"id"`
}
