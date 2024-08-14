package services

import (
	"errors"
	"log"

	"github.com/FredericoBento/HandGame/domain/player"
	"github.com/FredericoBento/HandGame/domain/player/memory"
	"github.com/google/uuid"
)

var (
	ErrPasswordDontMatch = errors.New("player password dont match")
)

type PlayerAuthConfiguration func(pas *PlayerAuthService) error

type PlayerAuthService struct {
	players player.PlayerRepository
}

func NewPlayerAuthService(cfgs ...PlayerAuthConfiguration) (*PlayerAuthService, error) {
	pas := &PlayerAuthService{}

	for _, cfg := range cfgs {
		err := cfg(pas)

		if err != nil {
			return nil, err
		}
	}
	return pas, nil
}

func WithPlayerRepository(pr player.PlayerRepository) PlayerAuthConfiguration {
	return func(pas *PlayerAuthService) error {
		pas.players = pr
		return nil
	}
}

func WithMemoryPlayerRepository() PlayerAuthConfiguration {
	pr := memory.New()
	return WithPlayerRepository(pr)
}

func (pas *PlayerAuthService) SignInPlayer(playerID uuid.UUID, password string) error {
	p, err := pas.players.Get(playerID)
	if err != nil {
		return err
	}

	success, err := pas.players.ComparePassword(p, password)
	if err != nil {
		return err
	}

	if success == false {
		return ErrPasswordDontMatch
	} else {
		log.Printf("Playuer has log-in successfuly")
	}

	return nil
}
