// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package config

type Snowplow struct {
	Enabled               bool   `json:"enabled"`
	StandardRoutesEnabled bool   `json:"standardRoutesEnabled"`
	OpenRedirectsEnabled  bool   `json:"openRedirectsEnabled"`
	GetPath               string `json:"getPath"`
	PostPath              string `json:"postPath"`
	RedirectPath          string `json:"redirectPath"`
}
