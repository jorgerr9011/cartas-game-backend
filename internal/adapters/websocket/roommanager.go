package websocket

import (
	"log"
	"sync"

	playerapp "github.com/jorgerr9011/cartas-game-backend/internal/app/player"
	roomapp "github.com/jorgerr9011/cartas-game-backend/internal/app/room"

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
	RoomID  room.RoomID
	Content []byte
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

			r.initializeGame(client)

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

			log.Printf("Mensaje desde el RoomManager: %v", message.Content)

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

		// termina bucle en tests
		case <-r.Stop:
			return
		}
	}
}

func (r *RoomManager) initializeGame(client *Client) {
	var player *player.Player
	var room *room.Room
	var err error

	// Aqui se busca si existe la sala
	if r.Rooms[client.RoomID] == nil {
		r.Rooms[client.RoomID] = make(map[*Client]bool)

		player, err = r.playerUseCase.CreatePlayer("nombreCliente")
		if err != nil {
			log.Printf("Error creando cliente: %v", err)
			return
		}
		log.Printf("Jugador creado: %#v", player)
	}

	if r.DomainRooms[client.RoomID] == nil {
		gamefactory := game.NewGameFactory()
		newGame := gamefactory.NewGame(client.GameName)

		room, err = r.roomUseCase.CreateRoom(client.RoomID, string(client.RoomID), newGame)
		if err != nil {
			log.Printf("Error creando sala: %v", err)
			return
		}
		log.Printf("Sala creada: %#v", room)
		log.Printf("Juego añadido: %#v", room.Game)

		r.DomainRooms[room.ID] = room
	} else {
		room = r.DomainRooms[client.RoomID]
		log.Printf("Sala ya existente: %#v", room)
	}

	if room == nil || player == nil {
		log.Printf("ERROR: room o player es nil antes de JoinRoom")
		return
	}

	if err := r.roomUseCase.JoinRoom(room.ID, player.ID); err != nil {
		log.Printf("Error uniendo jugador %v a la sala %v : %v", room.Name, player.Name, err)
	}

	r.Rooms[client.RoomID][client] = true
}
