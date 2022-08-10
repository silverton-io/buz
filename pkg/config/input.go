// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package config

type Inputs struct {
	Snowplow    `json:"snowplow"`
	Cloudevents `json:"cloudevents"`
	Generic     `json:"generic"`
	Webhook     `json:"webhook"`
	Pixel       `json:"pixel"`
	Relay       `json:"relay"`
}
