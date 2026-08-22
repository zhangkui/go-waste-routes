package redis

import (
	"context"
	"fmt"

	redisv9 "github.com/redis/go-redis/v9"
	"go-waste-routes/internal/platform/config"
)

func Open(ctx context.Context, cfg config.RedisConfig) (*redisv9.Client, error) {
	client := redisv9.NewClient(&redisv9.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}
