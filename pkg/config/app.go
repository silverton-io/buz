// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package config

type App struct {
	Version           string `json:"version"`
	Name              string `json:"name"`
	Env               string `json:"env"`
	Mode              string `json:"mode"`
	Port              string `json:"port"`
	TrackerDomain     string `json:"trackerDomain"`
	EnableConfigRoute bool   `json:"enableConfigRoute"`
}
