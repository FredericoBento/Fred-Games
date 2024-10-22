package ws

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/FredericoBento/HandGame/internal/utils"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Event    chan *Event
	Username string
	RoomCode string
}

type ReadMessageHandler func(*Client, []byte)

type ReadEventHandler func(*Client, []byte)

const (
	// writeWait      = 10 * time.Second
	// pongWait       = 15 * time.Second
	writeWait      = 10 * time.Second
	pongWait       = 1 * time.Second
	pingPeriod     = (pongWait * 9) / 20
	maxMessageSize = 512
)

func NewClient(conn *websocket.Conn, username string) *Client {
	return &Client{
		Conn:     conn,
		Event:    make(chan *Event),
		Username: username,
		RoomCode: "",
	}
}

func (client *Client) SendEvent(e *Event) {
	slog.Info("Sent event")
	client.Event <- e
}

func (client *Client) SendErrorEvent(e *Event) {
	e.IsError = true
	e.From = "server"
	e.To = client.Username
	e.Data = json.RawMessage{}
	e.RoomCode = client.RoomCode
	client.SendEvent(e)
}

func (client *Client) SendErrorEventWithMessage(e *Event, message string) {
	type Message struct {
		message string `json:"message"`
	}
	m, err := utils.EncodeJSON(Message{message: message})
	if err != nil {
		slog.Error("Could not send message in error event with message")
		return
	}
	e.IsError = true
	e.From = "server"
	e.To = client.Username
	e.Data = m
	e.RoomCode = client.RoomCode
	client.SendEvent(e)
}

func (client *Client) ReadPump(hub *Hub, handler ReadEventHandler) {
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, event, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Error: %v", err)
			}
			break
		}
		// Let service handle the messages
		handler(client, event)
	}
}

func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case event, ok := <-client.Event:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// if event.Data =
			eventBytes, err := utils.EncodeJSON(event)
			if err != nil {
				slog.Error("Error while marshiling: " + err.Error())
				return
			}
			w.Write(eventBytes)
			// slog.Info("Sent Message to " + client.Username + " with: " + string(message.Data) + " and " + string(message.Errors))
			// message.PrintMessage()
			// messageJson, err := json.Marshal(message)
			// if err != nil {
			// slog.Error("Error while marshiling: " + err.Error())
			// return
			// }
			// w.Write(messageJson)

			n := len(client.Event)
			for i := 0; i < n; i++ {
				e, ok := <-client.Event
				if !ok {
					client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				eventBytes, err := utils.EncodeJSON(e)
				if err != nil {
					slog.Error("Error while marshiling: " + err.Error())
					return
				}
				w.Write(eventBytes)

				if err := w.Close(); err != nil {
					return
				}
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
