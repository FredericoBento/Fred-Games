package services

import (
	"errors"
	"testing"

	"github.com/FredericoBento/HandGame/aggregate"
)

func TestPlayerAuth_NewPlayerAuthService(t *testing.T) {

	pas, err := NewPlayerAuthService(
		WithMemoryPlayerRepository(),
	)

	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		test           string
		expectedErr    error
		username       string
		password       string
		passwordToTest string
	}

	testCases := []testCase{
		{
			test:           "player password dont match",
			expectedErr:    ErrPasswordDontMatch,
			username:       "Fred",
			password:       "password123",
			passwordToTest: "password_is_different",
		},
		{
			test:           "player password match",
			expectedErr:    nil,
			username:       "Fred",
			password:       "password123",
			passwordToTest: "password123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			pl, err := aggregate.NewPlayer(tc.username, tc.password)
			if err != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("expected error: %v, got %v", tc.expectedErr, err)
				} else {
					t.Errorf("didnt expect to find this error here: %v", tc.expectedErr)
				}
			}

			err = pas.players.Add(pl)
			if err != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("expected error: %v, got %v", tc.expectedErr, err)
				} else {
					t.Errorf("didnt expect to find this error here: %v", tc.expectedErr)
				}
			}

			err = pas.SignInPlayer(pl.GetID(), tc.passwordToTest)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error: %v, got %v", tc.expectedErr, err)
			}
		})
	}

}
