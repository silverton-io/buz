package request

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMapParams(t *testing.T) {
	url := "something/else?p1=v1&p2=v2&p2=v3"
	want := map[string]string{
		"p1": "v1",
		"p2": "v2",
	}
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = req

	params := MapParams(c)

	equiv := reflect.DeepEqual(params, want)
	if !equiv {
		t.Fatalf(`got %v, want %v`, params, want)
	}
}
