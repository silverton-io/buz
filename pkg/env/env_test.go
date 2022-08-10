// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package env

import "testing"

func TestEnviron(t *testing.T) {
	want_dev := "development"
	want_stg := "staging"
	want_prd := "production"
	if DEV_ENVIRONMENT != want_dev {
		t.Fatalf(`DEV_ENVIRONMENT is %v, want %v`, DEV_ENVIRONMENT, want_dev)
	}
	if STG_ENVIRONMENT != want_stg {
		t.Fatalf(`STG_ENVIRONMENT is %v, want %v`, STG_ENVIRONMENT, want_stg)
	}
	if PROD_ENVIRONMENT != want_prd {
		t.Fatalf(`PROD_ENVIRONMENT is %v, want %v`, PROD_ENVIRONMENT, want_prd)
	}
}

func TestConfigPath(t *testing.T) {
	want := "HONEYPOT_CONFIG_PATH"
	if HONEYPOT_CONFIG_PATH != want {
		t.Fatalf(`HONEYPOT_CONFIG_PATH is %v, want %v`, HONEYPOT_CONFIG_PATH, want)
	}
}
