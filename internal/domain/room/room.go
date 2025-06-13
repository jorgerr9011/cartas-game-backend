package room

import (
	"errors"
	"time"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/game"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
)

type RoomID string

type Room struct {
	ID        RoomID
	Name      string
	Players   []player.PlayerID
	CreatedAt time.Time
	Started   bool
	TurnIndex int
	Game      game.Game
}

func NewRoom(id RoomID, name string, game game.Game) *Room {
	return &Room{
		ID:        id,
		Name:      name,
		Players:   []player.PlayerID{},
		CreatedAt: time.Now(),
		Started:   false,
		TurnIndex: -1,
		Game:      game,
	}
}

func (r *Room) AssignGame(g game.Game) {
	r.Game = g
}

func (r *Room) AddPlayer(playerID player.PlayerID) error {
	if r.Started {
		return errors.New("game already started")
	}
	for _, p := range r.Players {
		if p == playerID {
			return errors.New("player already in room")
		}
	}
	r.Players = append(r.Players, playerID)
	return nil
}

func (r *Room) StartGame() error {
	if r.Started {
		return errors.New("game already started")
	}
	if len(r.Players) < 2 {
		return errors.New("need at least two players to start")
	}
	r.Started = true
	r.TurnIndex = 0
	return nil
}

func (r *Room) NextTurn() {
	if !r.Started {
		return
	}
	r.TurnIndex = (r.TurnIndex + 1) % len(r.Players)
}

func (r *Room) CurrentPlayer() player.PlayerID {
	if r.TurnIndex == -1 || len(r.Players) == 0 {
		return ""
	}
	return r.Players[r.TurnIndex]
}
