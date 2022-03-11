package util

import "testing"

// TestMd5 calls Md5 with a string,
// checking to ensure it returns the appropriate hash.
func TestMd5(t *testing.T) {
	stringToHash := "giggitygiggitygoo"
	want := "c4f081a6f2bcd2d2a40441c161f46dca"
	h := Md5(stringToHash)
	if !(h == want) {
		t.Fatalf(`Md5(%q) = %q, want %q`, stringToHash, h, want)
	}
}

// TestMd5Empty calls Md5 with an empty string,
// checking to ensure it returns the appropriate hash.
func TestMd5Empty(t *testing.T) {
	empty := ""
	want := "d41d8cd98f00b204e9800998ecf8427e"
	h := Md5(empty)
	if !(h == want) {
		t.Fatalf(`Md5(%q) = %q, want %q`, empty, h, want)
	}
}
