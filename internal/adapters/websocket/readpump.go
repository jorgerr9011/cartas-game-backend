package websocket

import (
	"encoding/json"
	"log"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
)

func ReadPump(rm *RoomManager, client *Client) {
	defer func() {
		rm.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var jsonMsg Message
		if err := json.Unmarshal(message, &jsonMsg); err != nil {
			log.Printf("Error deserializando mensaje %s: %v", client.ID, err)
			return
		}

		if jsonMsg.RoomID == "" {
			jsonMsg.RoomID = client.RoomID // fallback si no se incluye el RoomID explícitamente
		}

		handleMessage(rm, client, jsonMsg)

		// rm.Broadcast <- Message{
		// 	Type:    jsonMsg.Type,
		// 	RoomID:  jsonMsg.RoomID,
		// 	Payload: jsonMsg.Payload,
		// }
	}
}

func handleMessage(rm *RoomManager, client *Client, msg Message) {
	switch msg.Type {
	case "start_game":
		// comienzo de juego
		room := rm.DomainRooms[msg.RoomID]
		room.Game.Start(room.Players)
		state := room.Game.GetState()
		log.Printf("Estado al empezar el juego: %#v", state)

		hand, err := room.Game.GetPlayerHand(player.PlayerID(client.Username + "-id"))
		if err != nil {
			log.Printf("Error obteniendo mano del jugador: %v", err)
		}
		log.Printf("Mano del jugador %v: %#v", client.Username, hand)

	case "player_move":
		//

	case "chat":
		//

	default:
		log.Printf("Tipo de mensaje desconocido: %s", msg.Type)
	}
}
