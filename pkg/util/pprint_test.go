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
