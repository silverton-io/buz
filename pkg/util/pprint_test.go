// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package util

import (
	"encoding/json"
	"testing"
)

func TestStringify(t *testing.T) {
	somethingToStringify := map[string]interface{}{
		"something": "here",
	}
	marshaled, err := json.Marshal(somethingToStringify)
	want := string(marshaled)
	got := Stringify(somethingToStringify)
	if got != want || err != nil {
		t.Fatalf(`Stringify(%v) = %v, want %v, err %v`, somethingToStringify, got, want, err)
	}
}
