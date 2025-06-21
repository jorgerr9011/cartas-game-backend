package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jorgerr9011/cartas-game-backend/internal/adapters/db"
	"github.com/jorgerr9011/cartas-game-backend/internal/adapters/redis"
	wsadapter "github.com/jorgerr9011/cartas-game-backend/internal/adapters/websocket"
	playerapp "github.com/jorgerr9011/cartas-game-backend/internal/app/player"
	roomapp "github.com/jorgerr9011/cartas-game-backend/internal/app/room"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/card"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
	"github.com/jorgerr9011/cartas-game-backend/pkg/config"
)

type Message struct {
	Type    string          `json:"type"`
	RoomID  room.RoomID     `json:"roomid"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type ClienteSimulado struct {
	Nombre   string
	Conn     *websocket.Conn
	Recibido chan string
}

// Simulación mínima para test
func startTestServer() (*wsadapter.RoomManager, string) {

	// Cargar las variables del .env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}

	redisDb := db.NewRedisClient(*cfg)

	gin.SetMode(gin.TestMode)

	// memRepo := memory.NewRoomRepo()
	// playerRepo := memory.NewPlayerRepo()
	// roomUseCase := roomapp.NewUseCase(memRepo)
	// playerUseCase := playerapp.NewUseCase(playerRepo)

	// Guardar en Redis
	redisRoomRepo := redis.NewRedisRoomRepo(redisDb)
	redisPlayerRepo := redis.NewRedisPlayerRepo(redisDb)
	roomUseCase := roomapp.NewUseCase(redisRoomRepo)
	playerUseCase := playerapp.NewUseCase(redisPlayerRepo)

	rm := wsadapter.NewRoomManager(roomUseCase, playerUseCase)
	go rm.Run() // arranca el loop del room manager

	router := gin.New()
	router.GET("/ws/:roomID", ServeWs(rm))

	server := httptest.NewServer(router)

	return rm, server.URL
}

func TestWebSocketEcho(t *testing.T) {
	roomName := "room-multi"
	rm, clients, wg := conectarClientes(t, roomName)

	// Cierre de conexiones
	defer func() {
		for _, c := range clients {
			_ = c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
			_ = c.Conn.Close()
		}
	}()

	// Espera que el cliente sea registrado y goroutines activas
	time.Sleep(100 * time.Millisecond)

	enviarMensajeComienzo(t, clients[0].Conn, roomName)

	time.Sleep(1000 * time.Millisecond)

	enviarMensajeJugarCarta(t, clients[0].Conn, roomName)

	time.Sleep(1000 * time.Millisecond)

	for _, c := range clients {
		leerTodosMensajes(t, c, 300*time.Millisecond)
	}

	// Detener goroutines
	done := make(chan struct{})
	close(done)

	// Espera a que goroutines terminen
	wg.Wait()

	// Limpieza
	close(rm.Stop)
}

func conectarClientes(t *testing.T, roomName string) (*wsadapter.RoomManager, []ClienteSimulado, *sync.WaitGroup) {
	rm, serverURL := startTestServer()
	u := strings.Replace(serverURL, "http", "ws", 1)

	var wg sync.WaitGroup
	done := make(chan struct{})
	var clients []ClienteSimulado

	// Conectar múltiples clientes
	usernames := []string{"jorge", "pepito", "ana", "luis"}
	for _, username := range usernames {
		url := fmt.Sprintf("%s/ws/%s?game_name=culo&username=%s", u, roomName, url.QueryEscape(username))
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("Error al conectar usuario %s: %v \n", username, err)
		}

		recibido := make(chan string, 20)
		iniciarReceptor(ws, recibido, done, &wg)

		clients = append(clients, ClienteSimulado{
			Nombre:   username,
			Conn:     ws,
			Recibido: recibido,
		})
	}

	return rm, clients, &wg
}

func enviarMensajeComienzo(t *testing.T, ws *websocket.Conn, roomName string) {
	msg := Message{
		Type:    "start_game",
		RoomID:  room.RoomID(roomName),
		Payload: json.RawMessage([]byte(`{}`)),
	}
	msgJSON, _ := json.Marshal(msg)
	if err := ws.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
		t.Fatalf("Error al enviar start_game desde cliente 1: %v \n", err)
	}
}

func enviarMensajeJugarCarta(t *testing.T, ws *websocket.Conn, roomName string) {
	player := "jorge-id"

	payload := wsadapter.PlayCardPayload{
		PlayerID: player,
		Card: card.Card{
			Suit: "bastos",
			Rank: 3,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error al serializar el payload: %v \n", err)
	}

	msg := Message{
		Type:    "play_card",
		RoomID:  room.RoomID(roomName),
		Payload: json.RawMessage(payloadBytes),
	}

	msgJSON, _ := json.Marshal(msg)
	if err := ws.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
		t.Fatalf("Error al enviar play_game desde el jugador %v: %v \n", player, err)
	}
}

func iniciarReceptor(ws *websocket.Conn, recibido chan string, done <-chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				_, msg, err := ws.ReadMessage()
				if err != nil {
					continue
				}
				recibido <- string(msg)
			}
		}
	}()
}

func leerTodosMensajes(t *testing.T, c ClienteSimulado, timeout time.Duration) {
	t.Logf("Mensajes recibidos por %s:", c.Nombre)
	for {
		select {
		case msg := <-c.Recibido:
			t.Logf("  %s", msg)
		case <-time.After(timeout):
			return
		}
	}
}
