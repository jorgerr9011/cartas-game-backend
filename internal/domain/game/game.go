// internal/domain/game/game.go
package game

import (
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/card"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
)

type GameID string

type PlayerState struct {
	ID    player.PlayerID
	Hand  []card.Card // Cartas en su mano
	Score int
}

type GameState struct {
	Turn            int
	CurrentPlayerID player.PlayerID
	Players         map[player.PlayerID]*PlayerState // Estado de cada jugador
	Started         bool
	Finished        bool
	DiscardPile     []card.Card
	Deck            []card.Card
}

type Game interface {
	GetName() string
	Start(playerIDs []player.PlayerID) error
	Play(playerID player.PlayerID, card card.Card) (GameState, error)
	GetPlayerHand(playerID player.PlayerID) (*PlayerState, error)
	GetState() GameState
	GetCurrentTurnPlayer() player.PlayerID
	IsFinished() bool
	MarshalState() ([]byte, error)
	UnmarshalState(data []byte) error
}
