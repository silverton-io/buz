// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package config

type GenericRootConfig struct {
	RootKey string `json:"rootKey"`
}

type GenericRootAndChildConfig struct {
	RootKey   string `json:"rootKey"`
	SchemaKey string `json:"schemaKey"`
	DataKey   string `json:"dataKey"`
}

type Generic struct {
	Enabled  bool                      `json:"enabled"`
	Path     string                    `json:"path"`
	Contexts GenericRootConfig         `json:"contexts"`
	Payload  GenericRootAndChildConfig `json:"payload"`
}
