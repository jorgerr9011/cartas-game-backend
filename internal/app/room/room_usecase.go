package roomapp

import (
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/card"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/game"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
	"github.com/jorgerr9011/cartas-game-backend/internal/ports/repository"
)

type UseCase interface {
	CreateRoom(id room.RoomID, name string, game game.Game) (*room.Room, error)
	JoinRoom(roomID room.RoomID, playerID player.PlayerID) error
	StartGame(roomID room.RoomID) (game.GameState, error)
	NextTurn(roomID room.RoomID) error
	CurrentPlayer(roomID room.RoomID) (player.PlayerID, error)
	Play(roomID room.RoomID, playerId player.PlayerID, card card.Card) (game.GameState, error)
}

type roomUseCase struct {
	repo repository.RoomRepository
}

func NewUseCase(r repository.RoomRepository) UseCase {
	return &roomUseCase{repo: r}
}

func (uc *roomUseCase) CreateRoom(id room.RoomID, name string, game game.Game) (*room.Room, error) {
	r := room.NewRoom(id, name, game)
	err := uc.repo.Save(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (uc *roomUseCase) JoinRoom(roomID room.RoomID, playerID player.PlayerID) error {
	r, err := uc.repo.FindByID(roomID)
	if err != nil {
		return err
	}
	err = r.AddPlayer(playerID)
	if err != nil {
		return err
	}
	return uc.repo.Save(r)
}

func (uc *roomUseCase) StartGame(roomID room.RoomID) (game.GameState, error) {
	r, err := uc.repo.FindByID(roomID)
	if err != nil {
		return game.GameState{}, err
	}
	err = r.Game.Start(r.Players)
	if err != nil {
		return game.GameState{}, err
	}
	uc.repo.Save(r)

	return r.Game.GetState(), nil
}

func (uc *roomUseCase) NextTurn(roomID room.RoomID) error {
	r, err := uc.repo.FindByID(roomID)
	if err != nil {
		return err
	}
	r.NextTurn()
	return uc.repo.Save(r)
}

func (uc *roomUseCase) CurrentPlayer(roomID room.RoomID) (player.PlayerID, error) {
	r, err := uc.repo.FindByID(roomID)
	if err != nil {
		return "", err
	}
	return r.CurrentPlayer(), nil
}

func (uc *roomUseCase) Play(roomID room.RoomID, playerId player.PlayerID, card card.Card) (game.GameState, error) {
	r, err := uc.repo.FindByID(roomID)
	if err != nil {
		return game.GameState{}, err
	}

	state, err := r.Game.Play(playerId, card)
	uc.repo.Save(r)
	return state, err
}
