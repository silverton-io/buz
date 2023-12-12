package annotator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSchemaMetadata(t *testing.T) {
	testData := []struct {
		name     string
		input    []byte
		expected schemaMetadata
	}{
		{
			name:  "Valid JSON",
			input: []byte(`{"self":{"vendor":"testVendor","namespace":"testNamespace","version":"testVersion"},"disableValidation":true}`),
			expected: schemaMetadata{
				Vendor:            "testVendor",
				Namespace:         "testNamespace",
				Version:           "testVersion",
				DisableValidation: true,
			},
		},
	}

	for _, tc := range testData {
		t.Run(tc.name, func(t *testing.T) {
			result := getSchemaMetadata(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
