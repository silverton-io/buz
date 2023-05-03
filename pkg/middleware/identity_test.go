// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package middleware

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/response"
	"github.com/stretchr/testify/assert"
)

func testHandler(c *gin.Context) {
	time.Sleep(3 * time.Millisecond)
	c.JSON(http.StatusOK, response.Ok)
}

const TEST_URL = "/somepath"

func buildTestServer(conf config.Identity) *httptest.Server {
	// Set up gin, router, middleware
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Identity(conf))
	r.GET(TEST_URL, testHandler)
	return httptest.NewServer(r)
}

func TestIdentityNoCookie(t *testing.T) {
	fallback := "some-id"
	cookie := config.IdentityCookie{
		Enabled:  true,
		Name:     "nuid",
		Secure:   true,
		TtlDays:  365,
		Domain:   "some.domain",
		Path:     "/",
		SameSite: "Lax",
	}
	noneCookie := config.IdentityCookie(cookie)
	noneCookie.SameSite = "None"
	strictCookie := config.IdentityCookie(cookie)
	strictCookie.SameSite = "Strict"

	noneConf := config.Identity{
		Cookie:   noneCookie,
		Fallback: fallback,
	}
	laxConf := config.Identity{
		Cookie:   cookie,
		Fallback: fallback,
	}
	strictConf := config.Identity{
		Cookie:   strictCookie,
		Fallback: fallback,
	}

	testCases := []struct {
		identityConf config.Identity
		wantSameSite http.SameSite
	}{
		{identityConf: noneConf, wantSameSite: http.SameSiteNoneMode},
		{identityConf: laxConf, wantSameSite: http.SameSiteLaxMode},
		{identityConf: strictConf, wantSameSite: http.SameSiteStrictMode},
	}

	for _, tc := range testCases {
		testServer := buildTestServer(tc.identityConf)
		// Make request
		resp, _ := http.Get(testServer.URL + TEST_URL)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf(`got status code %v, want %v`, resp.StatusCode, http.StatusOK)
		}
		defer resp.Body.Close()
		identityCookie := resp.Cookies()[0]
		assert.Equal(t, cookie.Name, identityCookie.Name)
		assert.Equal(t, cookie.Domain, identityCookie.Domain)
		assert.Equal(t, cookie.Path, identityCookie.Path)
		assert.Equal(t, cookie.Secure, identityCookie.Secure)
		assert.Equal(t, tc.wantSameSite, identityCookie.SameSite)
	}
}

func TestIdentityWithCookie(t *testing.T) {
	fallback := "some-id"
	someCookieVal := "some-cookie-val"
	cookie := config.IdentityCookie{
		Enabled:  true,
		Name:     "nuid",
		Secure:   true,
		TtlDays:  365,
		Domain:   "some.domain",
		Path:     "/",
		SameSite: "Lax",
	}
	conf := config.Identity{
		Cookie:   cookie,
		Fallback: fallback,
	}
	testServer := buildTestServer(conf)
	// Set up cookiejar
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	c := &http.Cookie{
		Name:   cookie.Name,
		Value:  someCookieVal,
		Domain: cookie.Domain,
		MaxAge: 300,
	}

	req, _ := http.NewRequest("GET", testServer.URL+TEST_URL, nil)
	req.AddCookie(c)
	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf(`got status code %v, want %v`, resp.StatusCode, http.StatusOK)
	}
	defer resp.Body.Close()
	identityCookie := resp.Cookies()[0]
	assert.Equal(t, cookie.Name, identityCookie.Name)
	assert.Equal(t, someCookieVal, identityCookie.Value)
}
