package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/stats"
)

func TestStatsHandler(t *testing.T) {
	u := uuid.New()
	now := time.Now().UTC()
	m := meta.CollectorMeta{
		Version:       "1.0.x",
		InstanceId:    u,
		StartTime:     now,
		TrackerDomain: "somewhere.net",
		CookieDomain:  "somewhere.io",
	}
	s := stats.ProtocolStats{}
	s.Build()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	handler := StatsHandler(&m, &s)

	handler(c)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf(`StatsHandler returned %d, want %d`, resp.StatusCode, http.StatusOK)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
	marshaledMeta, err := json.Marshal(m)
	if err != nil {
		t.Fatalf(`Could not marshal meta`)
	}
	equiv := reflect.DeepEqual(b, marshaledMeta)
	if !equiv {
		t.Fatalf(`StatsHandler returned %v, want %v`, b, marshaledMeta)
	}
}
