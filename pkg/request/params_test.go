package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMapParams(t *testing.T) {
	url := "something/else?p1=v1&p2=v2&p2=v3"
	want := map[string]interface{}{
		"p1": "v1",
		"p2": "v2",
	}
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = req

	params := MapParams(c)
	assert.Equal(t, params, want)
}
