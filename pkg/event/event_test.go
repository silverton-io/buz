package event

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayloadAsByte(t *testing.T) {
	raw := map[string]interface{}{
		"something": "here",
	}
	p := Payload{
		"something": "here",
	}
	b, err := p.AsByte()
	raw_bytes, _ := json.Marshal(raw)
	assert.Equal(t, raw_bytes, b)
	assert.Equal(t, nil, err)
}
