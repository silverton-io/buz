// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package middleware

import (
	"bytes"
	"encoding/json"
	"io"
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
	Body                     interface{}   `json:"body"`
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
		end := time.Now().UTC()
		duration := util.GetDuration(start, end)
		buf, _ := io.ReadAll(c.Request.Body)
		r1 := io.NopCloser(bytes.NewBuffer(buf))
		r2 := io.NopCloser(bytes.NewBuffer(buf))
		reqBody, err := io.ReadAll(r1)
		c.Request.Body = r2
		c.Next()
		if err != nil {
			log.Error().Err(err).Msg("could not read request body")
		}

		var b interface{}
		if string(reqBody) != "" {
			err = json.Unmarshal(reqBody, &b)

			if err != nil {
				log.Debug().Err(err).Interface("body", reqBody).Msg("could not unmarshal request body")
			}
		}

		r := request{
			ResponseCode:             c.Writer.Status(),
			RequestDuration:          duration,
			RequestDurationForHumans: duration.String(),
			ClientIp:                 getIp(c),
			RequestMethod:            c.Request.Method,
			RequestUri:               c.Request.RequestURI,
			Body:                     b,
		}
		log.Info().Interface("request", r).Msg("🟢")
	}
}
