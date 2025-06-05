package player

import (
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
	"github.com/jorgerr9011/cartas-game-backend/internal/ports/repository"
)

type UseCase interface {
	CreatePlayer(name string) (*player.Player, error)
	FindByID(id player.PlayerID) (*player.Player, error)
	List() ([]*player.Player, error)
}

type playerUseCase struct {
	repo repository.PlayerRepository
}

func NewUseCase(r repository.PlayerRepository) UseCase {
	return &playerUseCase{repo: r}
}

func (uc *playerUseCase) CreatePlayer(name string) (*player.Player, error) {
	id := player.PlayerID(name + "-id")
	p := player.NewPlayer(id, name)
	err := uc.repo.Save(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *playerUseCase) FindByID(id player.PlayerID) (*player.Player, error) {
	p, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *playerUseCase) List() ([]*player.Player, error) {
	pl, err := uc.repo.List()
	if err != nil {
		return nil, err
	}
	return pl, nil
}
