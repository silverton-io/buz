package config

type Privacy struct {
	Anonymize  `json:"anonymize"`
	DoNotTrack `json:"doNotTrack"`
}

type DoNotTrack struct {
	Enabled bool `json:"enabled"`
}

type Anonymize struct {
	DeviceIp        bool `json:"deviceIp"`
	DeviceUseragent bool `json:"deviceUseragent"`
	UserId          bool `json:"userId"`
}
