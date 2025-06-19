package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/game"
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

	// var rm room.Room
	// err = json.Unmarshal(data, &rm)

	gamefactory := game.NewGameFactory()
	rm, err := LoadRoomFromRedis(data, *gamefactory)

	return rm, err
}

func (r *RedisRoomRepo) Save(rm *room.Room) error {
	key := fmt.Sprintf("room:%s", rm.ID)
	// data, err := json.Marshal(rm)
	data, err := rm.MarshalForRedis()
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), key, data, 30*time.Minute).Err()
}

func (r *RedisRoomRepo) List() ([]*room.Room, error) {
	//
	return nil, nil
}

func LoadRoomFromRedis(data []byte, gameFactory game.GameFactory) (*room.Room, error) {
	var dto room.RoomDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, err
	}

	gameInstance := gameFactory.NewGame(dto.GameName)
	if err := gameInstance.UnmarshalState(dto.GameState); err != nil {
		return nil, err
	}

	room := &room.Room{
		ID:        room.RoomID(dto.ID),
		Name:      dto.Name,
		Players:   dto.Players,
		CreatedAt: dto.CreatedAt,
		Started:   dto.Started,
		TurnIndex: dto.TurnIndex,
		Game:      gameInstance,
	}

	return room, nil
}
