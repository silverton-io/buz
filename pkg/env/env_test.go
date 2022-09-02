// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package env

import "testing"

func TestConfigPath(t *testing.T) {
	want := "BUZ_CONFIG_PATH"
	if BUZ_CONFIG_PATH != want {
		t.Fatalf(`BUZ_CONFIG_PATH is %v, want %v`, BUZ_CONFIG_PATH, want)
	}
}
