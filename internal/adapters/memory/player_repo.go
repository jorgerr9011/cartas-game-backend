package memory

import (
	"errors"
	"sync"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
	"github.com/jorgerr9011/cartas-game-backend/internal/ports/repository"
)

type PlayerRepo struct {
	mu    sync.RWMutex
	store map[player.PlayerID]*player.Player
}

func NewPlayerRepo() repository.PlayerRepository {
	return &PlayerRepo{store: make(map[player.PlayerID]*player.Player)}
}

func (r *PlayerRepo) Save(player *player.Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[player.ID] = player
	return nil
}

func (r *PlayerRepo) FindByID(id player.PlayerID) (*player.Player, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if rm, ok := r.store[id]; ok {
		return rm, nil
	}
	return nil, errors.New("room not found")
}

func (r *PlayerRepo) List() ([]*player.Player, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*player.Player, 0, len(r.store))
	for _, v := range r.store {
		out = append(out, v)
	}
	return out, nil
}
