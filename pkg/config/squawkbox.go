// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package config

type Squawkbox struct {
	Enabled         bool   `json:"enabled"`
	CloudeventsPath string `json:"cloudeventsPath"`
	SnowplowPath    string `json:"snowplowPath"`
	GenericPath     string `json:"genericPath"`
}
