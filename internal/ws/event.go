package ws

import "encoding/json"

type EventType int

type Event struct {
	Type     EventType       `json:"type"`
	Data     json.RawMessage `json:"data,omitempty"`
	RoomCode string          `json:"roomCode",omitempty`
	From     string          `json:"from,omitempty"`
	To       string          `json:"to",omitempty`
	IsError  bool            `json:"isError",omitempty`
}

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
