package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/response"
)

func TestHealthcheckHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	HealthcheckHandler(c)

	resp := rec.Result()

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf(`HealthcheckHandler returned status code %v, want %v`, resp.StatusCode, http.StatusOK)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	marshaledB, _ := json.Marshal(response.Ok)
	equiv := reflect.DeepEqual(b, marshaledB)
	if !equiv {
		t.Fatalf(`HealthcheckHandler returned body %v, want %v`, b, marshaledB)
	}
}
