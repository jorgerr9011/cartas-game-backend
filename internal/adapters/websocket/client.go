package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
)

type Client struct {
	ID       string
	Conn     *websocket.Conn
	Send     chan []byte
	RoomID   room.RoomID
	GameName string
}
