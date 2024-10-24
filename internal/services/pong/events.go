package pong

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/utils"
	"github.com/FredericoBento/HandGame/internal/ws"
)

type EventMessage struct {
	Message string `json:"message"`
}

type EventCreateRoomData struct {
	Code string `json:"code"`
}

type EventRoomCreatedData struct {
	Code string `json:"code"`
}

type EventJoinRoomData struct {
	Code   string `json:"code"`
	Player string `json:"player"`
}

type EventPlayerJoinedRoomData struct {
	Code   string `json:"code"`
	Player string `json:"player"`
}

type EventJoinedRoomData struct {
	Code   string `json:"code"`
	Player string `json:"player"`
}

const (
	EventTypeGameSettings = 0
	EventTypeMessage      = 1

	EventTypeCreateRoom  = 21
	EventTypeCreatedRoom = 22

	EventTypeJoinRoom         = 23
	EventTypeJoinedRoom       = 24
	EventTypePlayerJoinedRoom = 25

	EventTypePaddleUpPressed = 31
	EventTypePaddleUpRelease = 32

	EventTypePaddleDownPressed = 33
	EventTypePaddleDownRelease = 34

	EventTypeBallShot = 35

	EventTypePlayerDisconnected = 4
)

var (
	ErrServerError = errors.New("Server couldnt process request")
)

func (s *PongService) HandleEventMessage(event *ws.Event, client *ws.Client) {
	message := EventMessage{}
	err := json.Unmarshal(event.Data, &message)
	if err != nil {
		slog.Error("Invalid data for EventMessage: " + err.Error())
		return
	}
	slog.Info("Got Message from " + client.Username + ": " + message.Message)
}

func (s *PongService) HandleEventCreateRoom(event *ws.Event, client *ws.Client) {
	code := utils.RandomString(4)
	_, exist := s.Hub.Rooms[code]
	for exist {
		code = utils.RandomString(4)
		_, exist = s.Hub.Rooms[code]
	}
	room := ws.NewRoom(code, 2)
	s.Hub.Rooms[code] = room

	err := room.AddClient(client)
	if err != nil {
		client.SendErrorEvent(event)
		return
	}
	createdRoomEvent := ws.NewEvent(EventTypeCreatedRoom, room.Code, "", client.Username)

	eventData := EventRoomCreatedData{
		Code: room.Code,
	}
	bytes, err := utils.EncodeJSON(eventData)
	if err != nil {
		client.SendErrorEventWithMessage(event, err.Error())
		return
	}
	createdRoomEvent.Data = bytes
	client.SendEvent(&createdRoomEvent)
}

func (s *PongService) HandleEventJoinRoom(event *ws.Event, client *ws.Client) {
	data := EventJoinRoomData{}
	err := json.Unmarshal(event.Data, &data)
	if err != nil {
		slog.Error("Invalid data for JoinRoomData")
		client.SendErrorEventWithMessage(event, ErrServerError.Error())
		return
	}
	if _, ok := s.Hub.Rooms[data.Code]; !ok {
		client.SendErrorEventWithMessage(event, ErrInvalidCode.Error())
		return
	}
	room := s.Hub.Rooms[data.Code]
	err = room.AddClient(client)
	if err != nil {
		client.SendErrorEventWithMessage(event, "Could not enter room")
		return
	}
	var otherClientUsername string
	for _, c := range room.Clients {
		if c.Username != client.Username {
			otherClientUsername = c.Username
			eventData := EventPlayerJoinedRoomData{
				Code:   data.Code,
				Player: client.Username,
			}
			slog.Info(eventData.Code)
			bytes, err := utils.EncodeJSON(eventData)
			if err != nil {
				c.SendErrorEventWithMessage(event, err.Error())
				return
			}
			event := ws.NewEvent(EventTypePlayerJoinedRoom, data.Code, "server", c.Username)
			event.Data = bytes
			c.SendEvent(&event)
		}
	}
	if otherClientUsername == "" {
		client.SendErrorEventWithMessage(event, "Room is empty")
		return
	}
	eventData := EventJoinedRoomData{
		Code:   room.Code,
		Player: otherClientUsername,
	}
	bytes, err := utils.EncodeJSON(eventData)
	if err != nil {
		client.SendErrorEventWithMessage(event, err.Error())
		return
	}
	joinedEvent := ws.NewSimpleEvent(EventTypeJoinedRoom, client.Username)
	joinedEvent.Data = bytes
	client.SendEvent(&joinedEvent)
}
