package memory

import (
	"fmt"
	"sync"

	"github.com/FredericoBento/HandGame/aggregate"
	"github.com/FredericoBento/HandGame/domain/player"
	"github.com/google/uuid"
)

type MemoryRepository struct {
	players map[uuid.UUID]aggregate.Player
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		players: make(map[uuid.UUID]aggregate.Player),
	}
}

func (mr *MemoryRepository) Get(id uuid.UUID) (aggregate.Player, error) {
	if player, ok := mr.players[id]; ok {
		return player, nil
	}
	return aggregate.Player{}, player.ErrPlayerNotFound
}

func (mr *MemoryRepository) Add(p aggregate.Player) error {
	if mr.players == nil {
		mr.Lock()
		mr.players = make(map[uuid.UUID]aggregate.Player)
		mr.Unlock()
	}

	if _, ok := mr.players[p.GetID()]; ok {
		return fmt.Errorf("player already exists :%w", player.ErrAddPlayer)
	}
	mr.Lock()
	mr.players[p.GetID()] = p
	mr.Unlock()

	return nil
}

func (mr *MemoryRepository) Update(p aggregate.Player) error {
	if _, ok := mr.players[p.GetID()]; ok {
		return fmt.Errorf("player does not exist: %w", player.ErrUpdatePlayer)
	}

	mr.Lock()
	mr.players[p.GetID()] = p
	mr.Unlock()

	return nil
}

func (mr *MemoryRepository) ComparePassword(p aggregate.Player, password string) (bool, error) {
	if _, ok := mr.players[p.GetID()]; ok {
		return false, fmt.Errorf("player does not exist: %w", player.ErrUpdatePlayer)
	}

	equalPasswords := false
	mr.Lock()
	if p.GetPassword() == password {
		equalPasswords = true
	}
	mr.Unlock()

	return equalPasswords, nil
}
