package player

import "errors"

type PlayerID string

type Player struct {
	ID    PlayerID
	Name  string
	Score int
}

func NewPlayer(id PlayerID, name string) *Player {
	return &Player{
		ID:   id,
		Name: name,
	}
}

func (p *Player) AddScore(points int) error {
	if points < 0 {
		return errors.New("points must be positive")
	}
	p.Score += points
	return nil
}
