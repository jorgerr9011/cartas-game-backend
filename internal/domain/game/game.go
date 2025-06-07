// internal/domain/game/game.go
package game

import "github.com/jorgerr9011/cartas-game-backend/internal/domain/player"

type GameID string

type Game interface {
	GetID() GameID
	GetName() string
	GetPlayers() []player.PlayerID
	GetCurrentTurnPlayer() player.PlayerID
	IsFinished() bool

	Start(players []player.PlayerID) error
	Play(playerID player.PlayerID, data map[string]interface{}) error
	GetState() map[string]interface{}
}
