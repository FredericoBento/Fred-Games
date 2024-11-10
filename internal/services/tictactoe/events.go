package tictactoe

import (
	"encoding/json"
	"errors"

	"github.com/FredericoBento/HandGame/internal/utils"
	"github.com/FredericoBento/HandGame/internal/ws"
)

const (
	EventTypeCreateGame        = 1
	EventTypeJoinGame          = 2
	EventTypeJoinedGame        = 22
	EventTypeOtherPlayerJoined = 23

	EventTypeMakePlay = 3

	EventTypePlayerSendMessage = 4

	EventTypePlayerDisconnected = 5
	EventTypePlayerReconnected  = 6
	EventTypeError              = 7

	EventTypeBoardCellUpdate = 8
	EventTypeStateUpdate     = 9

	EventTypeTie     = 10
	EventTypeVictory = 11
	EventTypeDefeat  = 12
)

var (
	ErrInternal       = errors.New("A internal server error has occured, try again later")
	ErrNotReconnected = errors.New("Invalid client to reconnect")

	ErrCouldNotPlay = errors.New("Could not make play")
	ErrCouldNotJoin = errors.New("Could not join game")

	ErrCellOcuppied = errors.New("This cell is already filled")
	ErrNotYourTurn  = errors.New("Not your turn")
)

func (s *TicTacToeService) HandleEventCreateGame(event *ws.Event, client *ws.Client) {
	code := s.generateUniqueCode(4)
	s.GameStates[code] = NewGameState(code)
	s.Hub.Rooms[code] = ws.NewRoom(code, 2)

	joinEvent := ws.NewEvent(EventTypeJoinGame, code)
	type Code struct {
		Code string `json:"code"`
	}
	bytes, err := utils.EncodeJSON(&Code{Code: code})
	if err != nil {
		delete(s.Hub.Rooms, code)
		delete(s.GameStates, code)
		s.Log.Error(err.Error())
		client.SendErrorEventWithMessage(event, ErrInternal.Error())
		return
	}

	joinEvent.Data = bytes
	s.HandleEventJoinGame(&joinEvent, client)
}

func (s *TicTacToeService) HandleEventJoinGame(event *ws.Event, client *ws.Client) {
	type EventData struct {
		Code string `json:"code"`
	}
	data := EventData{}
	err := json.Unmarshal(event.Data, &data)
	if err != nil {
		s.SendError(event, ErrInternal, client)
		return
	}
	room, exists := s.Hub.Rooms[data.Code]
	if !exists {
		s.SendError(event, ErrInvalidCode, client)
		return
	}

	state, exists := s.GameStates[data.Code]
	if !exists {
		s.Log.Error("This shouldnt happen though")
		s.SendError(event, ErrInvalidCode, client)
		return
	}

	err = room.AddClient(client)
	if err != nil {
		s.Log.Error(err.Error())
		s.SendError(event, ErrCouldNotJoin, client)
	}

	created := false
	if state.Player1 == nil {
		state.Player1 = NewPlayer(client.Username)
		created = true
	} else {
		if state.Player1.Username == client.Username {
			err = s.PlayerReconnect(state, client)
			if err != nil {
				s.Log.Error(err.Error())
				otherErr := room.RemoveClient(client)
				if otherErr != nil {
					s.Log.Error(otherErr.Error())
				}
				s.SendError(event, ErrCouldNotJoin, client)
			}
			return
		}
	}

	if !created {
		if state.Player2 == nil {
			state.Player2 = NewPlayer(client.Username)
			if state.Player1 != nil {
				ev := ws.NewSimpleEvent(EventTypeOtherPlayerJoined)
				ev.Data, err = utils.EncodeJSON(&state.Player2)
				if err != nil {
					s.Log.Error(err.Error())
					return
				}
				s.Hub.Clients[state.Player1.Username].SendEvent(&ev)
			}
		} else {
			if state.Player2.Username == client.Username {
				err = s.PlayerReconnect(state, client)
				if err != nil {
					otherErr := room.RemoveClient(client)
					if otherErr != nil {
						s.Log.Error(otherErr.Error())
					}
					s.Log.Error(err.Error())
					s.SendError(event, ErrCouldNotJoin, client)
				}
				return
			}
		}

	}

	joinedEvent := ws.NewEvent(EventTypeJoinedGame, room.Code)

	joinedEvent.Data, err = utils.EncodeJSON(&state)
	if err != nil {
		s.Log.Error(err.Error())
		otherErr := room.RemoveClient(client)
		if otherErr != nil {
			s.Log.Error(otherErr.Error())
		}
		s.SendError(event, ErrCouldNotJoin, client)
	}

	client.SendEvent(&joinedEvent)
}

func (s *TicTacToeService) HandleEventMakePlay(event *ws.Event, client *ws.Client) {
	type Play struct {
		Row int `json:"row"`
		Col int `json:"col"`
	}
	play := Play{}
	err := json.Unmarshal(event.Data, &play)
	if err != nil {
		s.Log.Error(err.Error())
		s.SendError(event, ErrCouldNotPlay, client)
		return
	}

	if client.RoomCode == "" {
		return
	}

	if _, ok := s.Hub.Rooms[client.RoomCode]; !ok {
		s.Log.Error(err.Error())
		s.SendError(event, ErrCouldNotPlay, client)
		return
	}

	state, ok := s.GameStates[client.RoomCode]
	if !ok {
		s.Log.Error(err.Error())
		s.SendError(event, ErrCouldNotPlay, client)
		return
	}

	playerID := 0

	if state.Turn%2 == 0 {
		if state.Player1 == nil {
			s.Log.Error(err.Error())
			s.SendError(event, ErrCouldNotPlay, client)
			return
		}
		if state.Player1.Username == client.Username {
			// Allow
			playerID = 1
			err = state.MakePlay(1, play.Row, play.Col)
			if err != nil {
				s.Log.Error("Error that shouldnt happen: " + err.Error())
				return
			}
		} else {
			s.SendError(event, ErrNotYourTurn, client)
			return
		}
	} else {
		if state.Player2 == nil {
			s.SendError(event, ErrCouldNotPlay, client)
			return
		}
		if state.Player2.Username == client.Username {
			// Alow
			playerID = 2
			err = state.MakePlay(2, play.Row, play.Col)
			if err != nil {
				s.Log.Error("Error that shouldnt happen: " + err.Error())
				return
			}
		} else {
			s.SendError(event, ErrNotYourTurn, client)
			return
		}
	}

	_, _, player_num := state.CheckWin()
	s.Log.Info("Check Win", "player_num", player_num, "status", state.Status, "winner", state.Winner)
	switch player_num {
	case 0:
		state.Turn += 1
		if state.Status == game_status_finished {
			go s.BroadCastGameTie(state, play.Row, play.Col, playerID)
			state.Restart(false)
		} else {
			s.BroadCastBoardCellUpdate(state, play.Row, play.Col, playerID)
		}
		break
	case 1:
		state.Player1.Wins += 1
		s.BroadCastGameFinish(state, play.Row, play.Col, playerID)
		state.Restart(false)
		break
	case 2:
		state.Player2.Wins += 1
		s.BroadCastGameFinish(state, play.Row, play.Col, playerID)
		state.Restart(false)
		break
	default:
		s.Log.Error("This shouldnt happen after checking win")
		break
	}

}

func (s *TicTacToeService) PlayerDisconnect(client *ws.Client) {
	s.Log.Info("Player disconnect handling extra logic here", "code", client.RoomCode)
	if client.RoomCode != "" {
		room, ok := s.Hub.Rooms[client.RoomCode]
		if ok {
			room.RemoveClient(client)
		}
		state, ok := s.GameStates[room.Code]
		if !ok {
			s.Log.Error("State not found")
			return
		}
		event := ws.NewEvent(EventTypePlayerDisconnected, room.Code)
		if state.Player1 != nil && state.Player1.Username == client.Username {
			state.Player1.Connected = false
			data, err := utils.EncodeJSON(state.Player1)
			if err != nil {
				s.Log.Error(err.Error())
				return
			}
			event.Data = data
			if state.Player2 != nil && state.Player2.Connected {
				if c, ok := s.Hub.Clients[state.Player2.Username]; ok {
					c.SendEvent(&event)
				}
			} else {
				delete(s.GameStates, state.Code)
			}
		} else {
			if state.Player2 != nil && state.Player2.Username == client.Username {
				state.Player2.Connected = false
				data, err := utils.EncodeJSON(state.Player2)
				if err != nil {
					s.Log.Error(err.Error())
					return
				}
				event.Data = data
				if state.Player1 != nil && state.Player1.Connected {
					if c, ok := s.Hub.Clients[state.Player1.Username]; ok {
						c.SendEvent(&event)
					}
				} else {
					delete(s.GameStates, state.Code)
				}
			}
		}
	}
	delete(s.Hub.Clients, client.Username)
	close(client.Event)
}

func (s *TicTacToeService) PlayerReconnect(state *GameState, client *ws.Client) error {
	if state.Player1.Username == client.Username {
		state.Player1.Connected = true
		stateEvent := ws.NewEvent(EventTypeStateUpdate, state.Code)
		data, err := utils.EncodeJSON(&state)
		if err != nil {
			s.Log.Error(err.Error())
			return err
		}
		stateEvent.Data = data
		go client.SendEvent(&stateEvent)

		if state.Player2 == nil {
			return nil
		}

		if otherClient, ok := s.Hub.Clients[state.Player2.Username]; ok {
			ev := ws.NewEvent(EventTypePlayerReconnected, state.Code)
			data, err = utils.EncodeJSON(state.Player1)
			if err != nil {
				s.Log.Error(err.Error())
				return err
			}
			ev.Data = data
			otherClient.SendEvent(&ev)
		}
		return nil
	} else {
		if state.Player2.Username == client.Username {
			state.Player2.Connected = true
			stateEvent := ws.NewEvent(EventTypeStateUpdate, state.Code)
			data, err := utils.EncodeJSON(&state)
			if err != nil {
				s.Log.Error(err.Error())
				return err
			}
			stateEvent.Data = data
			go client.SendEvent(&stateEvent)

			if otherClient, ok := s.Hub.Clients[state.Player1.Username]; ok {
				ev := ws.NewEvent(EventTypePlayerReconnected, state.Code)
				data, err = utils.EncodeJSON(state.Player2)
				if err != nil {
					s.Log.Error(err.Error())
					return err
				}
				ev.Data = data
				otherClient.SendEvent(&ev)
			}
			return nil
		}
	}
	return ErrNotReconnected

}

func (s *TicTacToeService) HandleEventChatMSG(event *ws.Event, client *ws.Client) {
	_, ok := s.Hub.Rooms[client.RoomCode]
	if !ok {
		s.SendError(event, errors.New("Something went wrong"), client)
		return
	}

	type Data struct {
		Message string `json:"message"`
		From    string `json:"from,omitempty"`
	}

	data := Data{}
	err := json.Unmarshal(event.Data, &data)
	if err != nil {
		s.SendError(event, errors.New("Something went wrong"), client)
		return
	}
	data.From = client.Username

	err = json.Unmarshal(event.Data, &data)
	if err != nil {
		s.SendError(event, errors.New("Something went wrong"), client)
		return
	}

	event.RoomCode = client.RoomCode
	s.Hub.Broadcast <- event
}

func (s *TicTacToeService) SendError(event *ws.Event, err error, client *ws.Client) {
	if client != nil {
		s.Log.Error(err.Error())
		client.SendErrorEventWithMessage(event, err.Error())
	}
}

func (s *TicTacToeService) generateUniqueCode(length int) string {
	code := utils.RandomString(length)
	for {
		if _, ok := s.GameStates[code]; ok {
			code = utils.RandomString(4)
		} else {
			return code
		}
	}
}

func (s *TicTacToeService) BroadCastState(state *GameState) {
	sendEvent := ws.NewEvent(EventTypeStateUpdate, state.Code)
	data, err := utils.EncodeJSON(&state)
	if err != nil {
		s.Log.Error(err.Error())
		return
	}
	sendEvent.Data = data
	for _, client := range s.Hub.Rooms[state.Code].Clients {
		client.SendEvent(&sendEvent)
	}
}

func (s *TicTacToeService) BroadCastBoardCellUpdate(state *GameState, row int, col int, value int) {
	sendEvent := ws.NewEvent(EventTypeBoardCellUpdate, state.Code)
	type Data struct {
		Row   int `json:"row"`
		Col   int `json:"col"`
		Value int `json:"value"`
	}
	data, err := utils.EncodeJSON(Data{Row: row, Col: col, Value: value})
	if err != nil {
		s.Log.Error(err.Error())
		return
	}
	sendEvent.Data = data
	s.BroadCastEvent(state.Code, &sendEvent)
}

func (s *TicTacToeService) BroadCastGameFinish(state *GameState, row int, col int, value int) {
	if state.Status == game_status_finished && state.Winner == 0 {
		// Tie
		go s.BroadCastGameTie(state, row, col, value)
		return
	}
	victoryEvent := ws.NewEvent(EventTypeVictory, state.Code)
	defeatEvent := ws.NewEvent(EventTypeDefeat, state.Code)

	type Data struct {
		Winner  int     `json:"winner"`
		Player1 *Player `json:"player1"`
		Player2 *Player `json:"player2"`
		Row     int     `json:"row"`
		Col     int     `json:"col"`
		Value   int     `json:"value"`
	}

	data, err := utils.EncodeJSON(Data{Winner: state.Winner, Player1: state.Player1,
		Player2: state.Player2, Row: row, Col: col, Value: value})

	if err != nil {
		s.Log.Error("Could not encode json when broadcasting game finished")
		return
	}

	victoryEvent.Data = data
	defeatEvent.Data = data
	s.Log.Info("state winner: ", state.Winner)
	if state.Winner == 1 {
		go s.Hub.Clients[state.Player1.Username].SendEvent(&victoryEvent)
		go s.Hub.Clients[state.Player2.Username].SendEvent(&defeatEvent)
	}
	if state.Winner == 2 {
		go s.Hub.Clients[state.Player2.Username].SendEvent(&victoryEvent)
		go s.Hub.Clients[state.Player1.Username].SendEvent(&defeatEvent)
	}
	if state.Winner == 0 {
		s.Log.Error("state status is game finish but is nto ")
	}
}

func (s *TicTacToeService) BroadCastGameTie(state *GameState, row int, col int, value int) {
	tieEvent := ws.NewEvent(EventTypeTie, state.Code)
	type Data struct {
		Tie   int `json:"ties"`
		Row   int `json:"row"`
		Col   int `json:"col"`
		Value int `json:"value"`
	}

	data, err := utils.EncodeJSON(Data{Tie: state.Ties, Row: row, Col: col, Value: value})

	if err != nil {
		s.Log.Error("Could not encode json when broadcasting tied game")
		return
	}
	tieEvent.Data = data
	go s.BroadCastEvent(tieEvent.RoomCode, &tieEvent)
}

func (s *TicTacToeService) BroadCastEvent(code string, ev *ws.Event) {
	if _, ok := s.Hub.Rooms[code]; !ok {
		s.Log.Error("Room does not exist")
		return
	}
	for _, client := range s.Hub.Rooms[code].Clients {
		go client.SendEvent(ev)
	}
}
