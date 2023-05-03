// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package testutil

import (
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
)

const URL = "/test"

func TestHandler(c *gin.Context) {
	time.Sleep(3 * time.Millisecond)
}

func BuildTestServer(handlerFuncs ...gin.HandlerFunc) *httptest.Server {
	// Set up gin, router, middleware
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Use provided middleware/handlerfuncs
	for _, hf := range handlerFuncs {
		r.Use(hf)
	}
	r.GET(URL, TestHandler)
	return httptest.NewServer(r)
}

func BuildRecordedEngine() (*httptest.ResponseRecorder, *gin.Context, *gin.Engine) {
	rec := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(rec)
	return rec, context, engine
}
