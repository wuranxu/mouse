package core

import (
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limit rate.Limit
	burst int
	*rate.Limiter
}

func NewRateLimiter(limit float64, burst int) *RateLimiter {
	r := rate.Limit(limit)
	limiter := rate.NewLimiter(r, burst)
	return &RateLimiter{
		limit: r, Limiter: limiter,
	}
}
