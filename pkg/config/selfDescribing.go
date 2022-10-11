// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package config

type SelfDescribingRootConfig struct {
	RootKey string `json:"rootKey"`
}

type SelfDescribingRootAndChildConfig struct {
	RootKey   string `json:"rootKey"`
	SchemaKey string `json:"schemaKey"`
	DataKey   string `json:"dataKey"`
}

type SelfDescribing struct {
	Enabled  bool                             `json:"enabled"`
	Path     string                           `json:"path"`
	Contexts SelfDescribingRootConfig         `json:"contexts"`
	Payload  SelfDescribingRootAndChildConfig `json:"payload"`
}
