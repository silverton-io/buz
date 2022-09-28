// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package config

type Inputs struct {
	Snowplow       `json:"snowplow"`
	Cloudevents    `json:"cloudevents"`
	SelfDescribing `json:"selfDescribing"`
	Webhook        `json:"webhook"`
	Pixel          `json:"pixel"`
}
