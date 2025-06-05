package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jorgerr9011/cartas-game-backend/internal/adapters/memory"
	playerapp "github.com/jorgerr9011/cartas-game-backend/internal/app/player"
	roomapp "github.com/jorgerr9011/cartas-game-backend/internal/app/room"
	handler "github.com/jorgerr9011/cartas-game-backend/internal/ports/http"
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

	roomHandler := handler.NewRoomHandler(createRoomUseCase)
	playerHandler := handler.NewPlayerHandler(createPlayerUseCase)

	api := r.Group("/api")
	roomHandler.Register(api)
	playerHandler.Register(api)

	log.Println("🚀  Server listening on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
