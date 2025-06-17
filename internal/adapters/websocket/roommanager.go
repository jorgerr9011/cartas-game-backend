package websocket

import (
	"encoding/json"
	"log"
	"sync"

	playerapp "github.com/jorgerr9011/cartas-game-backend/internal/app/player"
	roomapp "github.com/jorgerr9011/cartas-game-backend/internal/app/room"

	"github.com/jorgerr9011/cartas-game-backend/internal/domain/card"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/game"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/player"
	"github.com/jorgerr9011/cartas-game-backend/internal/domain/room"
)

type RoomManager struct {
	Rooms       map[room.RoomID]map[*Client]bool //conexiones websocket activas
	DomainRooms map[room.RoomID]*room.Room
	Broadcast   chan Message
	Register    chan *Client
	Unregister  chan *Client
	Stop        chan struct{}
	mu          sync.RWMutex

	roomUseCase   roomapp.UseCase
	playerUseCase playerapp.UseCase
}

type Message struct {
	Type    string          `json:"type"`
	RoomID  room.RoomID     `json:"roomid"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type PlayCardPayload struct {
	PlayerID string    `json:"player_id"`
	Card     card.Card `json:"card"`
}

func NewRoomManager(roomUC roomapp.UseCase, playerUC playerapp.UseCase) *RoomManager {
	return &RoomManager{
		Rooms:       make(map[room.RoomID]map[*Client]bool),
		DomainRooms: make(map[room.RoomID]*room.Room),
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Stop:        make(chan struct{}),

		roomUseCase:   roomUC,
		playerUseCase: playerUC,
	}
}

func (r *RoomManager) Run() {
	for {
		select {
		case client := <-r.Register:
			r.mu.Lock()
			r.initializeRoom(client)
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

			jsonMsg, err := json.Marshal(message)
			if err != nil {
				log.Printf("\nError deserializando mensaje: %v", err)
			}

			clients := r.Rooms[message.RoomID]
			for c := range clients {
				select {
				case c.Send <- jsonMsg:
				default:
					close(c.Send)
					delete(clients, c)
				}
			}
			r.mu.Unlock()

		// termina bucle en tests
		case <-r.Stop:
			return
		}
	}
}

func (r *RoomManager) initializeRoom(client *Client) {
	var player *player.Player
	var room *room.Room
	var err error

	r.handleClientInRoom(client)
	player = r.handleCreatePlayer(client)
	room = r.handleCreateRoomGame(client)

	if room == nil || player == nil {
		log.Printf("\nERROR: room o player es nil antes de JoinRoom")
		return
	}

	if err = r.roomUseCase.JoinRoom(room.ID, player.ID); err != nil {
		log.Printf("\nError uniendo jugador %v a la sala %v : %v", room.Name, player.Name, err)
	}

	log.Printf("\nSala al final del registro: %#v", room)
}

func (r *RoomManager) handleCreatePlayer(client *Client) *player.Player {
	player, err := r.playerUseCase.CreatePlayer(client.Username)
	if err != nil {
		log.Printf("\nError creando cliente: %v", err)
		return nil
	}
	log.Printf("\nJugador creado: %#v", player)
	return player
}

func (r *RoomManager) handleCreateRoomGame(client *Client) *room.Room {
	if r.DomainRooms[client.RoomID] == nil {
		gamefactory := game.NewGameFactory()
		newGame := gamefactory.NewGame(client.GameName)

		room, err := r.roomUseCase.CreateRoom(client.RoomID, string(client.RoomID), newGame)
		if err != nil {
			log.Printf("\nError creando sala: %v", err)
			return nil
		}
		log.Printf("\nSala creada: %#v", room)
		log.Printf("\nJuego añadido: %#v", room.Game)

		r.DomainRooms[room.ID] = room
	}
	room := r.DomainRooms[client.RoomID]
	log.Printf("\nSala ya existente: %#v", room)
	return room
}

func (r *RoomManager) handleClientInRoom(client *Client) {
	if r.Rooms[client.RoomID] == nil {
		r.Rooms[client.RoomID] = make(map[*Client]bool)
	}
	r.Rooms[client.RoomID][client] = true
}
