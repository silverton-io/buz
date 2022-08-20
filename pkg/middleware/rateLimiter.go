// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/response"
	limiter "github.com/ulule/limiter/v3"
	ginMiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func getDurationFromString(period string) time.Duration {
	switch period {
	case "MS":
		return 1 * time.Millisecond
	case "S":
		return 1 * time.Second
	case "M":
		return 1 * time.Minute
	case "H":
		return 1 * time.Hour
	case "D":
		return 24 * time.Hour
	default:
		return 1 * time.Second
	}
}

func onLimitReachedHandler(c *gin.Context) {
	log.Trace().Stack().Msg("limit reached - throttled request")
	c.JSON(http.StatusTooManyRequests, response.RateLimitExceeded)
}

func BuildRateLimiter(conf config.RateLimiter) *limiter.Limiter {
	period := getDurationFromString(conf.Period)
	rate := limiter.Rate{
		Period: period,
		Limit:  conf.Limit,
	}
	store := memory.NewStore()
	l := limiter.New(store, rate)
	return l
}

func BuildRateLimiterMiddleware(l *limiter.Limiter) gin.HandlerFunc {
	middleware := ginMiddleware.NewMiddleware(l, ginMiddleware.WithLimitReachedHandler(onLimitReachedHandler))
	return middleware
}
