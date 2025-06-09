package repository

import "github.com/jorgerr9011/cartas-game-backend/internal/domain/game"

type GameRepository interface {
	Save(g *game.CuloCardGame) error
	FindByID(id game.GameID) (*game.CuloCardGame, error)
	Delete(id game.GameID) error
	Update(g *game.CuloCardGame) error
}
