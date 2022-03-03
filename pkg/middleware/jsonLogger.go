package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/util"
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

func JsonAccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := util.GetDuration(start)
		r := request{
			ResponseCode:             c.Writer.Status(),
			RequestDuration:          duration,
			RequestDurationForHumans: duration.String(),
			ClientIp:                 getIp(c),
			RequestMethod:            c.Request.Method,
			RequestUri:               c.Request.RequestURI,
		}
		log.Info().Interface("request", r).Msg("")
	}
}
