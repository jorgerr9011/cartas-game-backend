package websocket

import "log"

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

		log.Printf("SALA DE JUEGO: %s", client.RoomID)
		log.Printf("OBJECT SALA: %#v", rm.DomainRooms[client.RoomID])

		rm.Broadcast <- Message{
			RoomID:  client.RoomID,
			Content: message,
		}
	}
}
