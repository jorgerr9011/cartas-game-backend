package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jorgerr9011/cartas-game-backend/internal/adapters/memory"
	"github.com/jorgerr9011/cartas-game-backend/internal/adapters/websocket"
	playerapp "github.com/jorgerr9011/cartas-game-backend/internal/app/player"
	roomapp "github.com/jorgerr9011/cartas-game-backend/internal/app/room"
	handlerHTTP "github.com/jorgerr9011/cartas-game-backend/internal/ports/http"
	"github.com/jorgerr9011/cartas-game-backend/internal/ports/ws"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	memRepo := memory.NewRoomRepo()
	playerRepo := memory.NewPlayerRepo()

	createRoomUseCase := roomapp.NewUseCase(memRepo)
	createPlayerUseCase := playerapp.NewUseCase(playerRepo)

	roomHandler := handlerHTTP.NewRoomHandler(createRoomUseCase)
	playerHandler := handlerHTTP.NewPlayerHandler(createPlayerUseCase)

	roommanager := websocket.NewRoomManager()
	go roommanager.Run()

	api := r.Group("/api")
	roomHandler.Register(api)
	playerHandler.Register(api)
	api.GET("/ws/:roomID", ws.ServeWs(roommanager))

	log.Println("🚀  Server listening on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
