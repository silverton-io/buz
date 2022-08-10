// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package util

import (
	"reflect"
	"testing"
)

type testStruct struct {
	A string `json:"a"`
}

// TestStructToMap ensures the proper map[string]interface{} is
// generated after calling it with an arbitrary struct.
func TestStructToMap(t *testing.T) {
	var s string = "something"
	ts := testStruct{
		A: s,
	}
	want := map[string]interface{}{
		"a": s,
	}
	got := StructToMap(ts)
	equivalent := reflect.DeepEqual(got, want)
	if !equivalent {
		t.Fatalf(`StructToMap(%v) = %v, want %v`, ts, got, want)
	}

}
