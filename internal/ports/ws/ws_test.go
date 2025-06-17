package ws

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jorgerr9011/cartas-game-backend/internal/adapters/memory"
	wsadapter "github.com/jorgerr9011/cartas-game-backend/internal/adapters/websocket"
	playerapp "github.com/jorgerr9011/cartas-game-backend/internal/app/player"
	roomapp "github.com/jorgerr9011/cartas-game-backend/internal/app/room"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/card"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
)

type Message struct {
	Type    string          `json:"type"`
	RoomID  room.RoomID     `json:"roomid"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// Simulación mínima para test
func startTestServer() (*wsadapter.RoomManager, string) {
	gin.SetMode(gin.TestMode)

	memRepo := memory.NewRoomRepo()
	playerRepo := memory.NewPlayerRepo()

	roomUseCase := roomapp.NewUseCase(memRepo)
	playerUseCase := playerapp.NewUseCase(playerRepo)

	rm := wsadapter.NewRoomManager(roomUseCase, playerUseCase)
	go rm.Run() // arranca el loop del room manager

	router := gin.New()
	router.GET("/ws/:roomID", ServeWs(rm))

	server := httptest.NewServer(router)

	// Devuelve el RoomManager y la URL del servidor para usarla en el test
	return rm, server.URL
}

func TestWebSocketEcho(t *testing.T) {
	rm, serverURL := startTestServer()

	u := strings.Replace(serverURL, "http", "ws", 1)
	roomName := "room-multi"

	// Conectar múltiples clientes
	var clients []*websocket.Conn
	usernames := []string{"jorge", "pepito", "ana", "luis"}
	for _, username := range usernames {
		url := fmt.Sprintf("%s/ws/%s?game_name=culo&username=%s", u, roomName, url.QueryEscape(username))
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("Error al conectar usuario %s: %v \n", username, err)
		}
		defer ws.Close()
		clients = append(clients, ws)
	}

	// Espera que el cliente sea registrado y goroutines activas
	time.Sleep(100 * time.Millisecond)

	// Comienzo del juego
	enviarMensajeComienzo(t, clients[0], roomName)

	// Jugar una carta
	enviarMensajeJugarCarta(t, clients[0], roomName)

	// Verificar que todos los clientes reciban un mensaje (broadcast)
	for i, c := range clients {
		_, data, err := c.ReadMessage()
		if err != nil {
			t.Errorf("Cliente %d: Error al leer mensaje: %v \n", i+1, err)
		} else {
			t.Logf("Cliente %d recibió: %s \n", i+1, string(data))
		}
	}

	// Solo validamos que el WebSocket esté vivo un instante
	t.Logf("WebSocket conectado correctamente a la room %s \n", "room-test")

	// Limpieza
	close(rm.Stop)
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
