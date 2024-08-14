package memory

import (
	"errors"
	"testing"

	"github.com/FredericoBento/HandGame/aggregate"
	"github.com/FredericoBento/HandGame/domain/player"
	"github.com/google/uuid"
)

func TestMemory_GetPlayer(t *testing.T) {
	type testCase struct {
		test        string
		id          uuid.UUID
		expectedErr error
	}

	pl, err := aggregate.NewPlayer("player1", "password123")
	if err != nil {
		t.Fatal(err)
	}

	id := pl.GetID()

	repo := MemoryRepository{
		players: map[uuid.UUID]aggregate.Player{
			id: pl,
		},
	}

	testCases := []testCase{
		{
			test:        "no player by id",
			id:          uuid.MustParse("26aa13e1-04af-4fe4-8e20-2e4b881ad6ae"),
			expectedErr: player.ErrPlayerNotFound,
		},
		{
			test:        "player by id",
			id:          id,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := repo.Get(tc.id)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemory_AddPlayer(t *testing.T) {
	type testCase struct {
		test        string
		expectedErr error
		player      aggregate.Player
	}

	player_exists, err := aggregate.NewPlayer("player1", "password123")
	if err != nil {
		t.Fatal(err)
	}

	player_dont_exists, err := aggregate.NewPlayer("player_dont_exits", "password123")
	if err != nil {
		t.Fatal(err)
	}

	id := player_exists.GetID()

	repo := MemoryRepository{
		players: map[uuid.UUID]aggregate.Player{
			id: player_exists,
		},
	}

	testCases := []testCase{
		{
			test:        "add player that do not exists",
			expectedErr: nil,
			player:      player_dont_exists,
		},
		{
			test:        "add player that already exists",
			expectedErr: player.ErrAddPlayer,
			player:      player_exists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := repo.Add(tc.player)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}

			_, err = repo.Get(tc.player.GetID())
			if err != nil {
				t.Errorf("Add played function gave no errors, but the repository has no player")
			}
		})
	}
}
