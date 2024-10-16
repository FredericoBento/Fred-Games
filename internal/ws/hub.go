package ws

import (
	"errors"
	"log/slog"
)

type Room struct {
	Code       string
	Clients    map[string]*Client
	MaxClients int
}

type Hub struct {
	Clients    map[string]*Client
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Event
}

var (
	ErrClientAlreadyInRoom = errors.New("Client is already in room")
	ErrRoomIsFull          = errors.New("Room Is full of clients")
)

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Event),
	}
}

func NewRoom(code string, maxClients int) *Room {
	return &Room{
		Code:       code,
		Clients:    make(map[string]*Client),
		MaxClients: maxClients,
	}
}

func (room *Room) AddClient(client *Client) error {
	if len(room.Clients) >= room.MaxClients {
		return ErrRoomIsFull
	}
	if _, ok := room.Clients[client.Username]; ok {
		return ErrClientAlreadyInRoom
	}
	room.Clients[client.Username] = client
	return nil
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client.Username] = client
			slog.Info("User " + client.Username + " has connected")
			// client.SendMessage(NewMessage([]byte("All Good"), "", ""))
			// client.SendSimpleMessage("All Good")
			// if _, ok := hub.Rooms[client.RoomCode]; ok {
			// room := hub.Rooms[client.RoomCode]

			// if _, ok := room.Clients[client.Username]; !ok {
			// room.Clients[client.Username] = client
			// }

			// }
		case client := <-hub.Unregister:
			slog.Info("User " + client.Username + " has disconnected")
			if _, ok := hub.Rooms[client.RoomCode]; ok {
				if _, ok := hub.Rooms[client.RoomCode].Clients[client.Username]; ok {
					if len(hub.Rooms[client.RoomCode].Clients) != 0 {
						// msg := NewMessage("player-left", &MessageData{Data: client.Username + " has left"}, nil)
						// msg := NewMessage("player-left", []byte(client.Username+" has left"), nil)
						// msg.RoomCode = client.RoomCode
						// msg.Username = client.Username
						// hub.Broadcast <- msg
					}
					delete(hub.Clients, client.Username)
					delete(hub.Rooms[client.RoomCode].Clients, client.Username)
					close(client.Event)
				}
			}
		case event := <-hub.Broadcast:
			if event.RoomCode != "" {
				if _, ok := hub.Rooms[event.RoomCode]; ok {
					for _, client := range hub.Rooms[event.RoomCode].Clients {
						client.Event <- event
					}
				}
			} else {
				if client, ok := hub.Clients[event.To]; ok {
					client.Event <- event
				}
			}

		}
	}
}
