package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestAdvancingCookieNoCookie(t *testing.T) {
	u := "/test"
	conf := config.Cookie{
		Enabled: true,
		Name:    "nomnomnom",
		Secure:  true,
		TtlDays: 3,
		Domain:  "sesame.street",
		Path:    "/",
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AdvancingCookie(conf))
	r.GET(u, testHandler)

	ts := httptest.NewServer(r)

	t.Run("no cookie", func(t *testing.T) {
		resp, _ := http.Get(ts.URL + u)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, resp.Header["Set-Cookie"])
		assert.Equal(t, "; Path=/; Domain=sesame.street; Max-Age=259200; Secure", resp.Header["Set-Cookie"][0][46:])
		assert.Equal(t, "nomnomnom=", resp.Header["Set-Cookie"][0][:10])
	})

	t.Run("pre-existing cookie", func(t *testing.T) {
		id := uuid.New()
		cookieVal := conf.Name + "=" + id.String()
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodGet, ts.URL+u, nil)
		req.Header.Set("Cookie", cookieVal)
		resp, _ := client.Do(req)
		wantSetCookie := "nomnomnom=" + id.String() + "; Path=/; Domain=sesame.street; Max-Age=259200; Secure"

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, resp.Header["Set-Cookie"])
		assert.Equal(t, []string{wantSetCookie}, resp.Header["Set-Cookie"])
	})
}
