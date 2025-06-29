package game

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/card"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
)

type CuloCardGame struct {
	// ID          GameID
	Players     []player.PlayerID
	TurnIndex   int
	TurnPlayer  int
	Started     bool
	Finished    bool
	Deck        []card.Card
	Played      []string
	PlayerHands map[player.PlayerID]*PlayerState
	DiscardPile []card.Card
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
	return g.Players[g.TurnPlayer]
}
func (g *CuloCardGame) IsFinished() bool { return g.Finished }

func (g *CuloCardGame) Start(playerIDs []player.PlayerID) error {

	hands, deck := g.barajarRepartirCartas(playerIDs)

	g.Players = playerIDs
	g.PlayerHands = hands
	g.Deck = deck
	g.DiscardPile = []card.Card{}
	g.TurnIndex = 1
	g.TurnPlayer = 0
	g.Started = true
	g.Finished = false

	return nil
}

func (g *CuloCardGame) Play(playerID player.PlayerID, card card.Card) (GameState, error) {
	playerIdFinded := g.GetCurrentTurnPlayer()
	if playerIdFinded != playerID {
		return GameState{}, fmt.Errorf("no es turno del jugador con ID: %s, jugador que le toca: %v", playerID, playerIdFinded)
	}

	/* Esto simplemente para hacer pruebas */
	cardFinded := g.PlayerHands[playerID].Hand[1]

	// g.jugarCarta(playerID, card)
	g.jugarCarta(playerID, cardFinded)

	g.avanzarTurno()

	return g.GetState(), nil
}

func (g *CuloCardGame) GetPlayerHand(playerID player.PlayerID) (*PlayerState, error) {
	hand, ok := g.PlayerHands[playerID]
	if !ok {
		return nil, fmt.Errorf("\nNo se encontró la mano del jugador %s", playerID)
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

func (g *CuloCardGame) avanzarTurno() {
	g.TurnIndex += 1
	g.TurnPlayer = (g.TurnPlayer + 1) % len(g.Players)
}

func (g *CuloCardGame) barajarRepartirCartas(playerIDs []player.PlayerID) (map[player.PlayerID]*PlayerState, []card.Card) {
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

	return hands, deck
}

func (g *CuloCardGame) jugarCarta(playerID player.PlayerID, card card.Card) error {
	manojugador := g.PlayerHands[playerID].Hand

	indiceCarta := g.findCardIndex(manojugador, card)
	if indiceCarta == -1 {
		return fmt.Errorf("\nNo se encontró la carta en la mano del jugador %s", playerID)
	}

	g.PlayerHands[playerID].Hand = slices.Delete(manojugador, indiceCarta, indiceCarta+1)

	g.DiscardPile = append(g.DiscardPile, card)

	return nil
}

func (g *CuloCardGame) findCardIndex(cards []card.Card, target card.Card) int {
	for i, c := range cards {
		if c.Suit == target.Suit && c.Rank == target.Rank {
			return i
		}
	}
	return -1
}

func (g *CuloCardGame) MarshalState() ([]byte, error) {
	currentPlayer := g.GetCurrentTurnPlayer()

	state := GameState{
		Turn:            g.TurnIndex,
		CurrentPlayerID: currentPlayer,
		Players:         g.PlayerHands,
		Started:         g.Started,
		Finished:        g.Finished,
		DiscardPile:     g.DiscardPile,
		Deck:            g.Deck,
	}

	return json.Marshal(state)
}

func (g *CuloCardGame) UnmarshalState(data []byte) error {
	var state GameState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	var playerIDs []player.PlayerID
	for id := range state.Players {
		playerIDs = append(playerIDs, id)
	}

	g.Players = playerIDs
	g.TurnIndex = state.Turn
	g.PlayerHands = state.Players
	g.Started = state.Started
	g.Finished = state.Finished
	g.DiscardPile = state.DiscardPile
	g.Deck = state.Deck

	g.TurnPlayer = -1
	for i, id := range g.Players {
		if id == state.CurrentPlayerID {
			g.TurnPlayer = i
			break
		}
	}

	return nil
}
