// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	A string `json:"a"`
}

// TestStructToMap ensures the proper map[string]interface{} is
// generated after calling it with an arbitrary struct.
func TestStructToMapGood(t *testing.T) {
	var s string = "something"
	ts := testStruct{
		A: s,
	}
	want := map[string]interface{}{
		"a": s,
	}
	got, _ := StructToMap(ts)
	equivalent := reflect.DeepEqual(got, want)
	if !equivalent {
		t.Fatalf(`StructToMap(%v) = %v, want %v`, ts, got, want)
	}

}

func TestStructToMapBad(t *testing.T) {
	something := "this"
	want := map[string]interface{}{}

	got, err := StructToMap(something)

	assert.Equal(t, want, got)
	assert.Error(t, err)
}
