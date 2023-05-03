package handler

import (
	"io"
	"net/http"
	"testing"

	testutil "github.com/silverton-io/buz/pkg/testUtil"
	"github.com/stretchr/testify/assert"
)

func TestBuzHandler(t *testing.T) {
	srv := testutil.BuildTestServer(BuzHandler())

	resp, _ := http.Get(srv.URL + testutil.URL)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf(`got status code %v, want %v`, resp.StatusCode, http.StatusOK)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "🐝", string(b))
}
