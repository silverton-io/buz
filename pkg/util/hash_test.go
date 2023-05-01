// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"fmt"
	"testing"
)

// TestMd5 calls Md5 with a string,
// checking to ensure it returns the appropriate hash.
func TestMd5(t *testing.T) {
	var tests = []struct {
		in   string
		want string
	}{
		{"giggitygiggitygoo", "c4f081a6f2bcd2d2a40441c161f46dca"},
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
	}
	for _, tt := range tests {
		tName := fmt.Sprintf("%v,%v", tt.in, tt.want)
		t.Run(tName, func(t *testing.T) {
			out := Md5(tt.in)
			if out != tt.want {
				t.Fatalf(`Md5(%v) = %v, want %v`, tt.in, out, tt.want)
			}
		})
	}
}
