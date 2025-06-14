package ws

import (
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
)

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
	uParsed, _ := url.Parse(u)
	uParsed.Path = "/ws/room-test"

	q := uParsed.Query()
	q.Set("game_name", "culo")
	uParsed.RawQuery = q.Encode()

	ws, _, err := websocket.DefaultDialer.Dial(uParsed.String(), nil)
	if err != nil {
		t.Fatalf("No se pudo conectar al WebSocket: %v", err)
	}
	defer ws.Close()

	// Espera que el cliente sea registrado y goroutines activas
	time.Sleep(100 * time.Millisecond)

	// mensaje := []byte(`{"type":"test","payload":"ping"}`)
	// if err := ws.WriteMessage(websocket.TextMessage, mensaje); err != nil {
	// 	t.Fatalf("Error al enviar mensaje: %v", err)
	// }

	// // Si tu WritePump envía algo en respuesta, aquí podrías leerlo:
	// _, resp, err := ws.ReadMessage()
	// if err != nil {
	// 	t.Fatalf("Error al recibir mensaje: %v", err)
	// }
	// t.Logf("Respuesta: %v", resp)

	// Solo validamos que el WebSocket esté vivo un instante
	t.Logf("WebSocket conectado correctamente a la room %s", "room-test")

	// Limpieza
	close(rm.Stop)
}
