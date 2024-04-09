package services

import (
	"net/http"
)

type RateLimitStatus = int

const (
	IP_REQUEST_LIMIT RateLimitStatus = iota + 1
	TOKEN_REQUEST_LIMIT
	REQUEST_ALLOWED
)

type RateLimiterChecker interface {
	CheckLimit(*http.Request) (RateLimitStatus, error)
}

type RateLimiter struct {
}

func (rateLimiter *RateLimiter) CheckLimit(*http.Request) (RateLimitStatus, error) {
	return IP_REQUEST_LIMIT, nil
}
