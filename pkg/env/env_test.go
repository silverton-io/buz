// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

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
	want := "BUZ_CONFIG_PATH"
	if BUZ_CONFIG_PATH != want {
		t.Fatalf(`BUZ_CONFIG_PATH is %v, want %v`, BUZ_CONFIG_PATH, want)
	}
}
