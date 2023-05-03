// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package middleware

import (
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/response"
	testutil "github.com/silverton-io/buz/pkg/testUtil"
	"github.com/stretchr/testify/assert"
)

func testHandler(c *gin.Context) {
	time.Sleep(3 * time.Millisecond)
	c.JSON(http.StatusOK, response.Ok)
}

const (
	TEST_COOKIE_FALLBACK = "some-cookie-fallback"
)

func TestIdentityNoCookie(t *testing.T) {
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
		Fallback: TEST_COOKIE_FALLBACK,
	}
	laxConf := config.Identity{
		Cookie:   cookie,
		Fallback: TEST_COOKIE_FALLBACK,
	}
	strictConf := config.Identity{
		Cookie:   strictCookie,
		Fallback: TEST_COOKIE_FALLBACK,
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
		testServer := testutil.BuildTestServer(Identity(tc.identityConf))
		// Make request
		resp, _ := http.Get(testServer.URL + testutil.URL)
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
		Fallback: TEST_COOKIE_FALLBACK,
	}
	testServer := testutil.BuildTestServer(Identity(conf))
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

	req, _ := http.NewRequest("GET", testServer.URL+testutil.URL, nil)
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
