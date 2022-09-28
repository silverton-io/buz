// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package config

type Config struct {
	App        `json:"app"`
	Middleware `json:"middleware"`
	Inputs     `json:"inputs"`
	Registry   `json:"registry"`
	Manifold   `json:"manifold,omitempty"`
	Sinks      []Sink `json:"sinks"`
	Squawkbox  `json:"squawkBox"`
	Privacy    `json:"privacy"`
	Tele       `json:"tele"`
}
