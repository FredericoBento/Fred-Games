package websocket_manager

import "golang.org/x/net/websocket"

type WebsocketManager struct {
	Users map[string]*User
}

type User struct {
	Username string
	Socket   *websocket.Conn
}

func (wm *WebsocketManager) AddUser(username string, socket *websocket.Conn) {
	wm.Users[username] = &User{
		Username: username,
		Socket:   socket,
	}
}

func (wm *WebsocketManager) RemoveUser(username string) {
	delete(wm.Users, username)
}

func (wm *WebsocketManager) GetUser(username string) *User {
	return wm.Users[username]
}
