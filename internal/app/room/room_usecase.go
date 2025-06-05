package roomapp

import (
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
	"github.com/jorgerr9011/cartas-game-backend/internal/ports/repository"
)

type UseCase interface {
	CreateRoom(name string) (*room.Room, error)
	JoinRoom(roomID room.RoomID, playerID room.PlayerID) error
	StartGame(roomID room.RoomID) error
	NextTurn(roomID room.RoomID) error
	CurrentPlayer(roomID room.RoomID) (room.PlayerID, error)
}

type roomUseCase struct {
	repo repository.RoomRepository
}

func NewUseCase(r repository.RoomRepository) UseCase {
	return &roomUseCase{repo: r}
}

func (uc *roomUseCase) CreateRoom(name string) (*room.Room, error) {
	id := room.RoomID(name + "-id")
	r := room.NewRoom(id, name)
	err := uc.repo.Save(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (uc *roomUseCase) JoinRoom(roomID room.RoomID, playerID room.PlayerID) error {
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

func (uc *roomUseCase) StartGame(roomID room.RoomID) error {
	r, err := uc.repo.FindByID(roomID)
	if err != nil {
		return err
	}
	err = r.StartGame()
	if err != nil {
		return err
	}
	return uc.repo.Save(r)
}

func (uc *roomUseCase) NextTurn(roomID room.RoomID) error {
	r, err := uc.repo.FindByID(roomID)
	if err != nil {
		return err
	}
	r.NextTurn()
	return uc.repo.Save(r)
}

func (uc *roomUseCase) CurrentPlayer(roomID room.RoomID) (room.PlayerID, error) {
	r, err := uc.repo.FindByID(roomID)
	if err != nil {
		return "", err
	}
	return r.CurrentPlayer(), nil
}
