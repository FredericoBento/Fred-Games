package tictactoe

import (
	"github.com/FredericoBento/HandGame/internal/ws"
)

func (s *TicTacToeService) Run(hub *ws.Hub) {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client.Username] = client
			s.Log.Info("User " + client.Username + " has connected")
			break

		case client := <-hub.Unregister:
			s.Log.Info("User " + client.Username + " has disconnected")

			if client.Conn.Close() != nil {
				s.Log.Error("Could not close connection")
			}

			s.PlayerDisconnect(client)
			break

		case event := <-hub.Broadcast:
			if event.RoomCode != "" {
				if _, ok := hub.Rooms[event.RoomCode]; ok {
					for _, client := range hub.Rooms[event.RoomCode].Clients {
						client.Event <- event
					}
				}
			}
			break
		}
	}
}
