package db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	return getClient(ctx)
}

func getClient(ctx context.Context) (*redis.Client, error) {
	fmt.Println(viper.GetString("REDIS_HOST"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_HOST"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(rdb)

	return rdb, nil
}
