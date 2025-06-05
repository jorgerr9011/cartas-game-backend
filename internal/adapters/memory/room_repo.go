package memory

import (
	"errors"
	"sync"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
	"github.com/jorgerr9011/cartas-game-backend/internal/ports/repository"
)

type RoomRepo struct {
	mu    sync.RWMutex
	store map[room.RoomID]*room.Room
}

func NewRoomRepo() repository.RoomRepository {
	return &RoomRepo{store: make(map[room.RoomID]*room.Room)}
}

func (r *RoomRepo) Save(room *room.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[room.ID] = room
	return nil
}

func (r *RoomRepo) FindByID(id room.RoomID) (*room.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if rm, ok := r.store[id]; ok {
		return rm, nil
	}
	return nil, errors.New("room not found")
}

func (r *RoomRepo) List() ([]*room.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*room.Room, 0, len(r.store))
	for _, v := range r.store {
		out = append(out, v)
	}
	return out, nil
}
