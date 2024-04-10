package services

import (
	"context"
	"errors"

	"github.com/nimbo1999/go-rate-limit/internal/repositories"
)

type RateLimitStatus = int

const (
	IP_REQUEST_LIMIT RateLimitStatus = iota + 1
	TOKEN_REQUEST_LIMIT
	REQUEST_ALLOWED
)

var (
	ErrIpNotFound     = errors.New("the ip was not provided")
	ErrApiKeyNotFound = errors.New("the API_KEY was not provided")
)

type RateLimiterChecker interface {
	CheckIpLimit(ip string) (RateLimitStatus, error)
	CheckApiKeyLimit(apiKey string) (RateLimitStatus, error)
}

type rateLimiterRedisChecker struct {
	repo                        repositories.RateLimiterRepository
	ipLimitCallsPerDuration     int
	apiKeyLimitCallsPerDuration int
}

func (rateLimiter *rateLimiterRedisChecker) CheckIpLimit(ip string) (RateLimitStatus, error) {
	if ip == "" {
		return IP_REQUEST_LIMIT, ErrIpNotFound
	}
	ipCalls, err := rateLimiter.repo.GetIpConnections(context.Background(), ip)
	if err != nil {
		return IP_REQUEST_LIMIT, err
	}
	if ipCalls >= rateLimiter.ipLimitCallsPerDuration {
		return IP_REQUEST_LIMIT, nil
	}
	go func(value string) {
		rateLimiter.repo.IncrementIpConnection(context.Background(), value)
	}(ip)
	return REQUEST_ALLOWED, nil
}

func (rateLimiter *rateLimiterRedisChecker) CheckApiKeyLimit(apiKey string) (RateLimitStatus, error) {
	if apiKey == "" {
		return TOKEN_REQUEST_LIMIT, ErrApiKeyNotFound
	}
	apiKeyCalls, err := rateLimiter.repo.GetApiKeyConnections(context.Background(), apiKey)
	if err != nil {
		return TOKEN_REQUEST_LIMIT, err
	}
	if apiKeyCalls >= rateLimiter.apiKeyLimitCallsPerDuration {
		return TOKEN_REQUEST_LIMIT, nil
	}
	go func(value string) {
		rateLimiter.repo.IncrementApiKeyConnection(context.Background(), value)
	}(apiKey)
	return REQUEST_ALLOWED, nil
}

func NewRateLimiterRedisChecker(
	repo repositories.RateLimiterRepository,
	ipLimitCallsPerDuration int,
	apiKeyLimitCallsPerDuration int,
) RateLimiterChecker {
	return &rateLimiterRedisChecker{
		repo:                        repo,
		ipLimitCallsPerDuration:     ipLimitCallsPerDuration,
		apiKeyLimitCallsPerDuration: apiKeyLimitCallsPerDuration,
	}
}
