package repository

import (
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
)

type PlayerRepository interface {
	Save(*player.Player) error
	FindByID(player.PlayerID) (*player.Player, error)
	List() ([]*player.Player, error)
}
