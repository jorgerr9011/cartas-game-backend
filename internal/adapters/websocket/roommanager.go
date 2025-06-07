package websocket

import "sync"

type RoomManager struct {
	Rooms      map[string]map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

type Message struct {
	RoomID  string
	Content []byte
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (r *RoomManager) Run() {
	for {
		select {
		case client := <-r.Register:
			r.mu.Lock()
			// Aqui se busca si existe la sala
			if r.Rooms[client.RoomID] == nil {
				r.Rooms[client.RoomID] = make(map[*Client]bool)
			}
			r.Rooms[client.RoomID][client] = true
			r.mu.Unlock()

		case client := <-r.Unregister:
			r.mu.Lock()
			if clients, ok := r.Rooms[client.RoomID]; ok {
				delete(clients, client)
				close(client.Send)
			}
			r.mu.Unlock()

		case message := <-r.Broadcast:
			r.mu.Lock()
			clients := r.Rooms[message.RoomID]
			for c := range clients {
				select {
				case c.Send <- message.Content:
				default:
					close(c.Send)
					delete(clients, c)
				}
			}
			r.mu.Unlock()
		}
	}
}
