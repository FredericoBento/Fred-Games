package ws

import (
	"encoding/json"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/utils"
)

type EventType int

type Event struct {
	Type     EventType       `json:"type"`
	Data     json.RawMessage `json:"data",omitempty`
	RoomCode string          `json:"roomCode",omitempty`
	From     string          `json:"from",omitempty"`
	To       string          `json:"to",omitempty`
	IsError  bool            `json:"isError",omitempty`
}

type EventPingPongData struct {
	Timestamp string `json:"timestamp"`
}

const (
	EventTypePing = 98
	EventTypePong = 99
)

func NewEvent(t EventType, roomCode string, from string, to string) Event {
	return Event{
		Type:     t,
		Data:     json.RawMessage{},
		RoomCode: roomCode,
		From:     from,
		To:       to,
		IsError:  false,
	}
}

func NewSimpleEvent(t EventType, to string) Event {
	return Event{
		Type:     t,
		Data:     json.RawMessage{},
		RoomCode: "",
		From:     "",
		To:       to,
		IsError:  false,
	}
}

func HandleEventPing(event *Event, client *Client) {
	data := EventPingPongData{}
	err := json.Unmarshal(event.Data, &data)
	if err != nil {
		slog.Error("could not pong")
		client.SendErrorEventWithMessage(event, "Error pinging")
		return
	}
	eventPong := NewSimpleEvent(EventTypePong, client.Username)
	data2 := EventPingPongData{
		Timestamp: data.Timestamp,
	}
	bytes, err := utils.EncodeJSON(data2)
	if err != nil {
		slog.Info("Error ping-pong")
		return
	}
	eventPong.Data = bytes
	client.SendEvent(&eventPong)
}
