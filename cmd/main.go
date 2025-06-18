package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jorgerr9011/cartas-game-backend/internal/adapters/db"
	"github.com/jorgerr9011/cartas-game-backend/internal/adapters/redis"
	"github.com/jorgerr9011/cartas-game-backend/internal/adapters/websocket"
	playerapp "github.com/jorgerr9011/cartas-game-backend/internal/app/player"
	roomapp "github.com/jorgerr9011/cartas-game-backend/internal/app/room"

	handlerHTTP "github.com/jorgerr9011/cartas-game-backend/internal/ports/http"
	"github.com/jorgerr9011/cartas-game-backend/internal/ports/ws"
	"github.com/jorgerr9011/cartas-game-backend/pkg/config"
)

func main() {
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		log.Fatal("No se pudo crear la carpeta logs:", err)
	}

	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("No se pudo abrir el archivo de log:", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	defer logFile.Close()

	// Cargar las variables del .env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}

	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redisDb := db.NewRedisClient(*cfg)

	r := gin.Default()

	// Guardar en memoria
	// memRepo := memory.NewRoomRepo()
	// playerRepo := memory.NewPlayerRepo()
	// roomUseCase := roomapp.NewUseCase(memRepo)
	// playerUseCase := playerapp.NewUseCase(playerRepo)

	// Guardar en Redis
	redisRoomRepo := redis.NewRedisRoomRepo(redisDb)
	redisPlayerRepo := redis.NewRedisPlayerRepo(redisDb)
	roomUseCase := roomapp.NewUseCase(redisRoomRepo)
	playerUseCase := playerapp.NewUseCase(redisPlayerRepo)

	roomHandler := handlerHTTP.NewRoomHandler(roomUseCase)
	playerHandler := handlerHTTP.NewPlayerHandler(playerUseCase)

	roommanager := websocket.NewRoomManager(roomUseCase, playerUseCase)
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
