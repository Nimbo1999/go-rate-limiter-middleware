package repositories

import "context"

type RateLimiterRepository interface {
	GetIpConnections(ctx context.Context, ip string) (int, error)
	IncrementIpConnection(ctx context.Context, ip string) error
	GetApiKeyConnections(ctx context.Context, apiKey string) (int, error)
	IncrementApiKeyConnection(ctx context.Context, apiKey string) error
}
