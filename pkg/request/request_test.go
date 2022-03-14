package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/silverton-io/honeypot/pkg/event"
)

func TestConst(t *testing.T) {
	want_js_ct := "application/json"

	if JSON_CONTENT_TYPE != want_js_ct {
		t.Fatalf(`got %v, want %v`, JSON_CONTENT_TYPE, want_js_ct)
	}
}

func TestPostEnvelope(t *testing.T) {
	u := "/somewhere"
	payload := event.SelfDescribingEnvelope{}
	marshaledPayload, _ := json.Marshal(payload)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("made it"))
		fmt.Println(r.URL.EscapedPath())
		d := r.URL.EscapedPath()
		if d != u {
			t.Fatalf(`posted payload to url %v, want %v`, d, u)
		}
		sentPayload, _ := ioutil.ReadAll(r.Body)
		payloadsEquiv := reflect.DeepEqual(sentPayload, marshaledPayload)
		if !payloadsEquiv {
			t.Fatalf(`posted body %v, want %v`, sentPayload, marshaledPayload)
		}
	}))
	defer ts.Close()

	dest, _ := url.Parse(ts.URL + u)

	PostEnvelope(*dest, payload)

}

func TestGet(t *testing.T) {}
