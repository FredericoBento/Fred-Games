package pong

type EventMessage struct {
	Message string `json:"message"`
}

type EventCreateRoomData struct {
	Code string `json:"code"`
}

type EventJoinRoomData struct {
	Code string `json:"code"`
}

const (
	EventTypeMessage    = 1
	EventTypeCreateRoom = 2
	EventTypeJoinRoom   = 3
	// EventTypePaddleMove         = 3
	// EventTypePlayerDisconnected = 4
)
