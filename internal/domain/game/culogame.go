package game

import (
	"errors"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
)

type CuloCardGame struct {
	// ID          GameID
	Players     []player.PlayerID
	TurnIndex   int
	Started     bool
	Finished    bool
	Deck        []string
	Played      []string
	PlayerHands map[player.PlayerID][]string
}

type GameState struct {
	Turn            int
	CurrentPlayerID player.PlayerID
	Players         []player.PlayerID
	Started         bool
	Finished        bool
}

func NewCuloCardGame() *CuloCardGame {
	return &CuloCardGame{
		// ID:          id,
		PlayerHands: make(map[player.PlayerID][]string),
	}
}

// func (g *CuloCardGame) GetID() GameID                 { return g.ID }
func (g *CuloCardGame) GetName() string               { return "CuloCardGame" }
func (g *CuloCardGame) GetPlayers() []player.PlayerID { return g.Players }
func (g *CuloCardGame) GetCurrentTurnPlayer() player.PlayerID {
	if len(g.Players) == 0 {
		return ""
	}
	return g.Players[g.TurnIndex]
}
func (g *CuloCardGame) IsFinished() bool { return g.Finished }

func (g *CuloCardGame) Start(players []player.PlayerID) error {
	if len(players) < 2 {
		return errors.New("not enough players")
	}
	g.Players = players
	g.TurnIndex = 0
	g.Started = true

	// Simulación de mazo
	g.Deck = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	// Distribuir 3 cartas por jugador
	for _, p := range players {
		g.PlayerHands[p] = g.Deck[:3]
		g.Deck = g.Deck[3:]
	}
	return nil
}

func (g *CuloCardGame) Play(playerID player.PlayerID, data map[string]interface{}) (GameState, error) {
	if g.GetCurrentTurnPlayer() != playerID {
		return g.GetState(), errors.New("not your turn")
	}
	card, ok := data["card"].(string)
	if !ok {
		return g.GetState(), errors.New("invalid card")
	}
	// Comprobar que el jugador tiene la carta
	hand := g.PlayerHands[playerID]
	index := -1
	for i, c := range hand {
		if c == card {
			index = i
			break
		}
	}
	if index == -1 {
		return g.GetState(), errors.New("card not in hand")
	}
	// Quitarla de la mano y jugarla
	g.PlayerHands[playerID] = append(hand[:index], hand[index+1:]...)
	g.Played = append(g.Played, card)

	// Avanzar turno
	g.TurnIndex = (g.TurnIndex + 1) % len(g.Players)

	// Finalizar si todos jugaron
	if len(g.Deck) == 0 {
		g.Finished = true
	}
	return g.GetState(), nil
}

func (g *CuloCardGame) GetState() GameState {
	return GameState{
		Turn:            g.TurnIndex,
		CurrentPlayerID: g.GetCurrentTurnPlayer(),
		Players:         g.Players,
		Started:         g.Started,
		Finished:        g.Finished,
	}
}
