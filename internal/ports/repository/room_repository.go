package repository

import "github.com/jorgerr9011/cartas-game-backend/internal/domain/room"

type RoomRepository interface {
	Save(*room.Room) error
	FindByID(room.RoomID) (*room.Room, error)
	List() ([]*room.Room, error)
}
