package memory

import (
	"errors"
	"sync"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/game"
	"github.com/jorgerr9011/cartas-game-backend/internal/ports/repository"
)

type GameRepo struct {
	mu    sync.RWMutex
	store map[game.GameID]*game.CuloCardGame
}

func NewGameRepo() repository.GameRepository {
	return &GameRepo{store: make(map[game.GameID]*game.CuloCardGame)}
}

func (r *GameRepo) Save(game *game.CuloCardGame) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[game.ID] = game
	return nil
}

func (r *GameRepo) FindByID(id game.GameID) (*game.CuloCardGame, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if rm, ok := r.store[id]; ok {
		return rm, nil
	}
	return nil, errors.New("room not found")
}

func (r *GameRepo) List() ([]*game.CuloCardGame, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*game.CuloCardGame, 0, len(r.store))
	for _, v := range r.store {
		out = append(out, v)
	}
	return out, nil
}

func (r *GameRepo) Update(g *game.CuloCardGame) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[g.ID] = g
	return nil
}

func (r *GameRepo) Delete(id game.GameID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.store, id)
	return nil
}
