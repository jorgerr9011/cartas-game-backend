package db

import (
	"context"

	"github.com/jorgerr9011/cartas-game-backend/pkg/config"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient(cfg config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Db_redis,
		Password: "",
		DB:       0,
	})
}
