package redis

import (
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
	"github.com/jorgerr9011/cartas-game-backend/internal/ports/repository"
	"github.com/redis/go-redis/v9"
)

type RedisPlayerRepo struct {
	client *redis.Client
}

func NewRedisPlayerRepo(client *redis.Client) repository.PlayerRepository {
	return &RedisPlayerRepo{client: client}
}

func (r *RedisPlayerRepo) Save(*player.Player) error {
	return nil
}

func (r *RedisPlayerRepo) FindByID(player.PlayerID) (*player.Player, error) {
	return nil, nil
}

func (r *RedisPlayerRepo) List() ([]*player.Player, error) {
	return nil, nil
}
