package pong

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
	Player string `json"player"`
}

type EventPlayerJoinedRoomData struct {
	Code   string `json:"code"`
	Player string `json"player"`
}

type EventJoinedRoomData struct {
	Code   string `json:"code"`
	Player string `json"player"`
}

const (
	EventTypeMessage      = 1
	EventTypeGameSettings = 11

	EventTypeCreateRoom  = 21
	EventTypeCreatedRoom = 22
	EventTypeJoinRoom    = 23

	EventTypeJoinedRoom       = 24
	EventTypePlayerJoinedRoom = 25

	EventTypePaddleUpPressed = 31
	EventTypePaddleUpRelease = 32

	EventTypePaddleDownPressed = 33
	EventTypePaddleDownRelease = 34

	EventTypeBallShot = 35

	EventTypePlayerDisconnected = 4
)
