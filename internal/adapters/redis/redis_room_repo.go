package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
	"github.com/jorgerr9011/cartas-game-backend/internal/ports/repository"
	"github.com/redis/go-redis/v9"
)

type RedisRoomRepo struct {
	client *redis.Client
}

func NewRedisRoomRepo(client *redis.Client) repository.RoomRepository {
	return &RedisRoomRepo{client: client}
}

func (r *RedisRoomRepo) FindByID(id room.RoomID) (*room.Room, error) {
	key := fmt.Sprintf("room:%s", id)
	data, err := r.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}
	var rm room.Room
	err = json.Unmarshal(data, &rm)
	return &rm, err
}

func (r *RedisRoomRepo) Save(rm *room.Room) error {
	key := fmt.Sprintf("room:%s", rm.ID)
	data, err := json.Marshal(rm)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), key, data, 30*time.Minute).Err()
}

func (r *RedisRoomRepo) List() ([]*room.Room, error) {
	//
	return nil, nil
}
