package websocket

import (
	"encoding/json"
	"log"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/game"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
)

func ReadPump(rm *RoomManager, client *Client) {
	defer func() {
		rm.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v \n", err)
			break
		}

		var jsonMsg Message
		if err := json.Unmarshal(message, &jsonMsg); err != nil {
			log.Printf("Error deserializando mensaje %s: %v \n", client.ID, err)
			return
		}

		if jsonMsg.RoomID == "" {
			jsonMsg.RoomID = client.RoomID // fallback si no se incluye el RoomID explícitamente
		}

		handleMessage(rm, client, jsonMsg)
	}
}

func handleMessage(rm *RoomManager, client *Client, msg Message) {
	switch msg.Type {
	case "start_game":
		// comienzo de juego

		// Esta obtención de room se haría en caso de guardar en memoria
		// room := rm.DomainRooms[msg.RoomID]

		state, err := rm.roomUseCase.StartGame(msg.RoomID)
		if err != nil {
			log.Printf("Error al inicial el juego %v.", err)
		}

		broadcastEstadoPartida(rm, "game_started", state, msg.RoomID)
		log.Printf("Juego empezado! Estado del juego: %#v", state)

	case "play_card":
		var payload PlayCardPayload
		err := json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			log.Fatalf("Error al deserializar payload: %v", err)
		}

		log.Printf("\nEl jugador %v jugará la carta %#v", payload.PlayerID, payload.Card)

		state, err := rm.roomUseCase.Play(msg.RoomID, player.PlayerID(payload.PlayerID), payload.Card)
		if err != nil {
			log.Printf("Error haciendo la jugada por parte del jugador %v.", err)
		}

		log.Printf("\n JUGADORES: %#v", rm.Rooms[msg.RoomID])

		log.Printf("\nEstado después de jugar una carta: %v", state)

		broadcastEstadoPartida(rm, "card_played", state, msg.RoomID)

		// room := rm.DomainRooms[msg.RoomID]

	case "end_game":
		//

	case "chat":
		//

	default:
		log.Printf("Tipo de mensaje desconocido: %s", msg.Type)
	}
}

func broadcastEstadoPartida(rm *RoomManager, tipo string, state game.GameState, roomID room.RoomID) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		log.Printf("Error comenzando el juego: %v", err)
	}

	log.Printf("Broadcast de tipo %s en room %s: %s", tipo, roomID, string(stateBytes))

	rm.Broadcast <- Message{
		Type:    tipo,
		RoomID:  roomID,
		Payload: json.RawMessage(stateBytes),
	}
}
