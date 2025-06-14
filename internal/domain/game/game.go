// internal/domain/game/game.go
package game

import "github.com/jorgerr9011/cartas-game-backend/internal/domain/player"

type GameID string

type Game interface {
	Start(players []player.PlayerID) error
	Play(playerID player.PlayerID, data map[string]interface{}) (GameState, error)
	GetState() GameState
	IsFinished() bool
}
