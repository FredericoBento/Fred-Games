package services

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/FredericoBento/HandGame/internal/mock"
	"github.com/FredericoBento/HandGame/internal/models"
)

func TestUserExists(t *testing.T) {
	t.Run("UserExists", func(t *testing.T) {
		mockRepo := &mock.MockUserRepository{
			GetByUsernameResult: &models.User{Username: "existing_user"},
			GetByUsernameError:  nil,
		}
		logger := slog.New(slog.Default().Handler())
		userService := &UserService{
			repo: mockRepo,
			log:  logger,
		}

		exists, err := userService.UserExists("existing_user")
		if !exists {
			t.Errorf("expected user to exist, got %v", exists)
		}
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

	})

	t.Run("UserDoesNotExists", func(t *testing.T) {
		expectedError := errors.New("user not found")
		mockRepo := &mock.MockUserRepository{
			GetByUsernameResult: nil,
			GetByUsernameError:  expectedError,
		}

		logger := slog.New(slog.NewJSONHandler(mock.MockIoWriter{}, &slog.HandlerOptions{}))
		userService := &UserService{
			repo: mockRepo,
			log:  logger,
		}

		exists, err := userService.UserExists("non_existing_user")
		if exists {
			t.Errorf("expected user to not exist, got %v", exists)
		}
		if err == nil || err != expectedError {
			t.Errorf("expected error '"+expectedError.Error()+"' , got %v", err)
		}

	})
}
