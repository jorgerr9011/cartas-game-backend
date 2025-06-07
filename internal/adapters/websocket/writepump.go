package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

func WritePump(client *Client) {
	defer client.Conn.Close()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				// Canal cerrado
				_ = client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Write error: %v", err)
				return
			}
		}
	}
}
