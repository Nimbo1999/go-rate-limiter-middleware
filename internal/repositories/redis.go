package repositories

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/nimbo1999/go-rate-limit/internal/utils"
	"github.com/redis/go-redis/v9"
)

func NewRateLimiterRedisRepository(client *redis.Client, expiration time.Duration) *rateLimiterRedisRepository {
	return &rateLimiterRedisRepository{client, expiration}
}

type rateLimiterRedisRepository struct {
	client     *redis.Client
	expiration time.Duration
}

func (repo *rateLimiterRedisRepository) GetIpConnections(ctx context.Context, ip string) (int, error) {
	return repo.GetDBValue(ctx, ip)
}

func (repo *rateLimiterRedisRepository) IncrementIpConnection(ctx context.Context, ip string) error {
	return repo.IncrementValue(ctx, ip)
}

func (repo *rateLimiterRedisRepository) GetApiKeyConnections(ctx context.Context, apiKey string) (int, error) {
	return repo.GetDBValue(ctx, apiKey)
}

func (repo *rateLimiterRedisRepository) IncrementApiKeyConnection(ctx context.Context, apiKey string) error {
	return repo.IncrementValue(ctx, apiKey)
}

func (repo *rateLimiterRedisRepository) IncrementValue(ctx context.Context, value string) error {
	intValue, err := repo.GetDBValue(ctx, value)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	intValue = intValue + 1
	return repo.client.Set(ctx, value, intValue, repo.expiration).Err()
}

func (repo *rateLimiterRedisRepository) GetDBValue(ctx context.Context, value string) (int, error) {
	ctx, cancel := utils.ContextWithMinimumTimeout(ctx)
	defer cancel()
	result, err := repo.client.Get(ctx, value).Result()
	if err != nil {
		if err == redis.Nil {
			// When redis.Nil it means that the key doesn't exist yet
			return 0, nil
		}
		fmt.Println(err.Error())
		return 0, err
	}
	if result == "" {
		return 0, nil
	}
	return strconv.Atoi(result)
}
