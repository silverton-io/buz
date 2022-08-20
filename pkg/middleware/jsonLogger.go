// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/util"
)

type request struct {
	ResponseCode             int           `json:"responseCode"`
	RequestDuration          time.Duration `json:"requestDuration"`
	RequestDurationForHumans string        `json:"requestDurationForHumans"`
	ClientIp                 string        `json:"clientIp"`
	RequestMethod            string        `json:"requestMethod"`
	RequestUri               string        `json:"requestUri"`
}

func getIp(c *gin.Context) string {
	ip := c.Request.Header.Get("X-Forwarded-For")
	if len(ip) == 0 {
		ip = c.Request.Header.Get("X-Real-IP")
	}
	if len(ip) == 0 {
		ip = c.Request.RemoteAddr
	}
	if strings.Contains(ip, ",") {
		ip = strings.Split(ip, ",")[0]
	}
	return ip
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		c.Next()
		end := time.Now().UTC()
		duration := util.GetDuration(start, end)
		r := request{
			ResponseCode:             c.Writer.Status(),
			RequestDuration:          duration,
			RequestDurationForHumans: duration.String(),
			ClientIp:                 getIp(c),
			RequestMethod:            c.Request.Method,
			RequestUri:               c.Request.RequestURI,
		}
		log.Info().Interface("🟢 request", r).Msg("")
	}
}
