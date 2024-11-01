package pong

import (
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/utils"
	"github.com/FredericoBento/HandGame/internal/ws"
)

func (s *PongService) Run(hub *ws.Hub) {
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
						s.GameStates[room.Code].RemovePlayer(client.Username)
						// hub.RemoveClientBroadcast(client)
						for _, c := range room.Clients {
							event := ws.NewSimpleEvent(ws.EventTypeUserDisconnected, "")
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
