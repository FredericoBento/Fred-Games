package aggregate

import (
	"errors"
	"testing"
)

func TestPlayer_NewPlayer(t *testing.T) {
	type testCase struct {
		test        string
		username    string
		password    string
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Empty username validation",
			username:    "",
			password:    "123",
			expectedErr: .ErrInvalidUser,
		}, {
			test:        "Valid username",
			username:    "Fred",
			password:    "123",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := aggregate.NewPlayer(tc.username, tc.password)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
