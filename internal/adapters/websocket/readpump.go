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

		// Esta obtención de room se haría en caso de guardar en memoria
		// room := rm.DomainRooms[msg.RoomID]

		state, err := rm.roomUseCase.StartGame(msg.RoomID)
		if err != nil {
			log.Printf("Error al inicial el juego %v.", err)
		}

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
		log.Printf("\nEstado después de jugar una carta: %v", state)

		// room := rm.DomainRooms[msg.RoomID]

		// turnPlayer := room.Game.GetCurrentTurnPlayer()
		// if turnPlayer != player.PlayerID(payload.PlayerID) {
		// 	log.Printf("El jugador que está intentando jugar no es al que le toca jugar este turno!")
		// }

		// // CARTA DE PRUEBA
		// playerHand, err := room.Game.GetPlayerHand(turnPlayer)
		// if err != nil {
		// 	log.Printf("Error obteniendo la mano del jugador")
		// }
		// log.Printf("Carta de prueba que se jugará %#v", playerHand.Hand[1])

		// // state, err := room.Game.Play(turnPlayer, payload.Card)
		// state, err := room.Game.Play(turnPlayer, playerHand.Hand[1])

	case "chat":
		//

	default:
		log.Printf("Tipo de mensaje desconocido: %s", msg.Type)
	}
}
