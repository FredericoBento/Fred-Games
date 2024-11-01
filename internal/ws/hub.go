package ws

import (
	"errors"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/utils"
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

const (
	EventTypeUserDisconnected = 4
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
	client.RoomCode = room.Code
	return nil
}

func (room *Room) RemoveClient(client *Client) error {
	if len(room.Clients) <= 0 {
		err := errors.New("Removing client of empty room")
		slog.Error(err.Error())
		return err
	}
	if _, ok := room.Clients[client.Username]; !ok {
		err := errors.New("Removing client from wrong room")
		slog.Error(err.Error())
		return err
	}
	client.RoomCode = ""
	delete(room.Clients, client.Username)
	return nil
}

func (hub *Hub) RemoveClientBroadcast(client *Client) error {
	event := NewSimpleEvent(EventTypeUserDisconnected, "")
	event.RoomCode = client.RoomCode
	type UsernameData struct {
		Username string `json:"username"`
	}
	data := UsernameData{
		Username: client.Username,
	}
	bytes, err := utils.EncodeJSON(data)
	if err != nil {
		slog.Error("Could not encode json while broadcasting cliennt disconnect")
		return err
	}
	event.Data = bytes
	hub.Broadcast <- &event
	slog.Info("broadcast here3")
	return nil
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client.Username] = client
			slog.Info("User " + client.Username + " has connected")
			if _, ok := hub.Rooms[client.RoomCode]; ok {
				room := hub.Rooms[client.RoomCode]

				if _, ok := room.Clients[client.Username]; !ok {
					room.Clients[client.Username] = client
				}
			}
		case client := <-hub.Unregister:
			slog.Info("User " + client.Username + " has disconnected")
			if _, ok := hub.Rooms[client.RoomCode]; ok {
				if _, ok := hub.Rooms[client.RoomCode].Clients[client.Username]; ok {
					if len(hub.Rooms[client.RoomCode].Clients) != 0 {
						room := hub.Rooms[client.RoomCode]
						err := room.RemoveClient(client)
						if err != nil {
							slog.Error("CRITICAL ERROR WHEN REMOVING CLIENT")
							return
						}
						delete(hub.Rooms[room.Code].Clients, client.Username)
						// hub.RemoveClientBroadcast(client)
						for _, c := range room.Clients {
							event := NewSimpleEvent(EventTypeUserDisconnected, "")
							event.RoomCode = client.RoomCode
							type UsernameData struct {
								Username string `json:"username"`
							}
							data := UsernameData{
								Username: client.Username,
							}
							bytes, err := utils.EncodeJSON(data)
							if err != nil {
								slog.Error("Could not encode json while broadcasting cliennt disconnect")
							} else {
								event.Data = bytes
								c.SendEvent(&event)
							}
						}
					}
				}
			}
			if client.Conn.Close() != nil {
				slog.Error("Could not close connection")
			}
			delete(hub.Clients, client.Username)
			close(client.Event)

		case event := <-hub.Broadcast:
			slog.Info("broadcast here4")
			if event.RoomCode != "" {
				if _, ok := hub.Rooms[event.RoomCode]; ok {
					for _, client := range hub.Rooms[event.RoomCode].Clients {
						event.To = client.Username
						slog.Info("Sent event to " + client.Username)
						client.Event <- event
					}
				}
			} else {
				if client, ok := hub.Clients[event.To]; ok {
					slog.Info("Sent event to " + client.Username)
					client.Event <- event
				}
			}

		}
	}
}
