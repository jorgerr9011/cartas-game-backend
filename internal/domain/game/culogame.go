package game

import (
	"fmt"
	"math/rand/v2"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/card"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
)

type CuloCardGame struct {
	// ID          GameID
	Players     []player.PlayerID
	TurnIndex   int
	Started     bool
	Finished    bool
	Deck        []card.Card
	Played      []string
	PlayerHands map[player.PlayerID]*PlayerState
	DiscardPile []card.Card
}

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

func NewCuloCardGame() *CuloCardGame {
	return &CuloCardGame{
		// ID:          id,
		// Players:     playerIDs,
		TurnIndex: -1,
		Started:   false,
		Finished:  false,
		// Deck:        deck,
		DiscardPile: []card.Card{},
		// PlayerHands: hands,
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

func (g *CuloCardGame) Start(playerIDs []player.PlayerID) error {
	deck := card.NewSpanishDeck()

	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	hands := make(map[player.PlayerID]*PlayerState)
	cardsPerPlayer := len(deck) / len(playerIDs)
	for _, playerid := range playerIDs {
		hand := deck[:cardsPerPlayer]
		deck = deck[cardsPerPlayer:]
		hands[playerid] = &PlayerState{
			ID:   playerid,
			Hand: hand,
		}
	}

	g.Players = playerIDs
	g.PlayerHands = hands
	g.Deck = deck
	g.DiscardPile = []card.Card{}
	g.TurnIndex = 0
	g.Started = true
	g.Finished = false

	return nil
}

func (g *CuloCardGame) Play(playerID player.PlayerID, data map[string]interface{}) (GameState, error) {
	//
	return g.GetState(), nil
}

func (g *CuloCardGame) GetPlayerHand(playerID player.PlayerID) (*PlayerState, error) {
	hand, ok := g.PlayerHands[playerID]
	if !ok {
		return nil, fmt.Errorf("no se encontró la mano del jugador %s", playerID)
	}
	return hand, nil
}

func (g *CuloCardGame) GetState() GameState {
	return GameState{
		Turn:            g.TurnIndex,
		CurrentPlayerID: g.GetCurrentTurnPlayer(),
		Players:         g.PlayerHands,
		Started:         g.Started,
		Finished:        g.Finished,
		DiscardPile:     g.DiscardPile,
		Deck:            g.Deck,
	}
}
