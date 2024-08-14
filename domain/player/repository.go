package player

import (
	"errors"

	"github.com/FredericoBento/HandGame/aggregate"
	"github.com/google/uuid"
)

var (
	ErrPlayerNotFound = errors.New("the player was not found in the repository")
	ErrAddPlayer      = errors.New("failed to add the player")
	ErrUpdatePlayer   = errors.New("failed to update the player")
)

type PlayerRepository interface {
	Get(uuid.UUID) (aggregate.Player, error)
	Add(aggregate.Player) error
	Update(aggregate.Player) error
	ComparePassword(aggregate.Player, string) (bool, error)
}
