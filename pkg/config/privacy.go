// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

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
