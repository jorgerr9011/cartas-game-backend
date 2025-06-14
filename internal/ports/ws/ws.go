package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	ws "github.com/jorgerr9011/cartas-game-backend/internal/adapters/websocket"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
)

// Revisar esto para producción
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(rm *ws.RoomManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := room.RoomID(c.Param("roomID"))
		gameName := c.Query("game_name")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.Copy().AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "WebSocket upgrade failed"})
		}
		defer conn.Close()

		client := &ws.Client{
			ID:       uuid.NewString(),
			Conn:     conn,
			Send:     make(chan []byte, 256),
			RoomID:   roomID,
			GameName: gameName,
		}

		rm.Register <- client

		go ws.WritePump(client)

		// Se ejecuta en la goroutine actual para que cuando termine ReadPump también termine la goroutine
		ws.ReadPump(rm, client)
	}
}
