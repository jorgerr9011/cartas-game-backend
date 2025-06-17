// internal/domain/game/game.go
package game

import (
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/card"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
)

type GameID string

type Game interface {
	Start(playerIDs []player.PlayerID) error
	Play(playerID player.PlayerID, card card.Card) (GameState, error)
	GetPlayerHand(playerID player.PlayerID) (*PlayerState, error)
	GetState() GameState
	GetCurrentTurnPlayer() player.PlayerID
	IsFinished() bool
}
