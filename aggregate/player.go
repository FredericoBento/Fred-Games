package aggregate

import (
	"errors"

	"github.com/FredericoBento/HandGame/entity"
	"github.com/google/uuid"
)

var (
	ErrInvalidUser = errors.New("a player has to have a valid username")
)

type Player struct {
	user *entity.User
}

func NewPlayer(username string, password string) (Player, error) {
	if username == "" || password == "" {
		return Player{}, ErrInvalidUser
	}

	user := &entity.User{
		ID:       uuid.New(),
		Username: username,
		Password: password,
	}

	return Player{
		user: user,
	}, nil
}

func (p *Player) GetID() uuid.UUID {
	return p.user.ID
}

func (p *Player) SetID(id uuid.UUID) {
	if p.user == nil {
		p.user = &entity.User{}
	}

	p.user.ID = id
}

func (p *Player) GetUsername() string {
	return p.user.Username
}

func (p *Player) SetUsername(username string) {
	if p.user == nil {
		p.user = &entity.User{}
	}

	p.user.Username = username
}

func (p *Player) GetPassword() string {
	return p.user.Password
}

func (p *Player) SetPassword(pwd string) {
	if p.user == nil {
		p.user = &entity.User{}
	}

	p.user.Password = pwd
}
