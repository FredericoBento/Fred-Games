package services

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/FredericoBento/HandGame/internal/mock"
	"github.com/FredericoBento/HandGame/internal/models"
)

func TestUserExists(t *testing.T) {
	t.Run("UserExists", func(t *testing.T) {
		mockRepo := &mock.MockUserRepository{
			GetByUsernameResult: &models.User{Username: "existing_user"},
			GetByUsernameError:  nil,
		}

		userService := NewUserService(mockRepo, 2*time.Minute)
		// userService.ChangeLogger(&slog.Logger{})

		exists, err := userService.UserExists(context.TODO(), "existing_user")
		if !exists {
			t.Errorf("expected user to exist, got %v", exists)
		}
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

	})

	t.Run("UserDoesNotExists", func(t *testing.T) {
		mockRepo := &mock.MockUserRepository{
			GetByUsernameResult: nil,
			GetByUsernameError:  sql.ErrNoRows,
		}

		userService := NewUserService(mockRepo, 2*time.Minute)
		// userService.ChangeLogger(slog.Default())

		exists, err := userService.UserExists(context.TODO(), "non_existing_user")
		if exists {
			t.Errorf("expected user to not exist, got %v", exists)
		}

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

	})
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		exists      bool
		existsErr   error
		createErr   error
		expectedErr error
	}{
		{
			name:        "User already exists",
			exists:      true,
			expectedErr: ErrUserAlreadyExists,
		},
		{
			name:        "Failed to check if user exist",
			existsErr:   errors.New("some error"),
			expectedErr: ErrUserExistsFailed,
		},
		{
			name:        "Successfully create user",
			exists:      false,
			existsErr:   sql.ErrNoRows,
			createErr:   nil,
			expectedErr: nil,
		},
		{
			name:        "Failed to create user",
			exists:      false,
			existsErr:   sql.ErrNoRows,
			createErr:   errors.New("repository failed to create user"),
			expectedErr: ErrCouldNotCreateUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockUserRepository{
				CreateError: tt.createErr,
			}

			user := &models.User{Username: "abc", Password: "abc"}

			if tt.exists {
				mockRepo.GetByUsernameError = nil
				mockRepo.GetByUsernameResult = user
			} else {
				mockRepo.GetByUsernameError = tt.existsErr
				mockRepo.GetByUsernameResult = nil
				mockRepo.CreateError = tt.createErr
			}

			us := NewUserService(mockRepo, 2*time.Minute)
			// us.ChangeLogger(slog.New(nil))

			err := us.CreateUser(context.Background(), user)

			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

		})
	}
}

func TestComparePassword(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		hashedPassword string
		expectedResult bool
		expectedErr    error
	}{
		{
			name:           "Password do match",
			password:       "abc",
			hashedPassword: "abc",
			expectedResult: true,
			expectedErr:    nil,
		},
		{
			name:           "Password do not match",
			password:       "abc",
			hashedPassword: "cba",
			expectedResult: false,
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockUserRepository{}

			s := NewUserService(mockRepo, 2*time.Minute)

			equal, err := s.ComparePassword(tt.password, tt.hashedPassword)

			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if equal != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, equal)
			}
		})
	}
}

func TestChangeLogger(t *testing.T) {
	tests := []struct {
		name        string
		logger      *slog.Logger
		expectedErr error
	}{
		{
			name:        "Invalid Logger",
			logger:      nil,
			expectedErr: ErrInvalidLogger,
		},
		{
			name:        "Valid Logger",
			logger:      slog.Default(),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockUserRepository{}

			s := NewUserService(mockRepo, 2*time.Minute)

			err := s.ChangeLogger(tt.logger)

			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	someUsers := []models.User{
		{
			Username: "abc",
			Password: "abc",
		},
		{
			Username: "cbc",
			Password: "cbc",
		},
	}

	tests := []struct {
		name           string
		expectedResult []models.User
		expectedErr    error
	}{
		{
			name:           "Failed to Get All",
			expectedResult: nil,
			expectedErr:    ErrCouldNotGetAllUsers,
		},
		{
			name:           "Successfuly Get All",
			expectedResult: someUsers,
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.MockUserRepository{
				GetAllResult: tt.expectedResult,
				GetAllError:  nil,
			}

			if tt.expectedResult == nil {
				mockRepo.GetAllError = errors.New("repository could not GetAll()")
				mockRepo.GetAllResult = nil
			}

			s := NewUserService(&mockRepo, 2*time.Minute)

			users, err := s.GetAllUsers()

			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			for i, user := range users {
				if user != tt.expectedResult[i] {
					t.Errorf("expected user %v, got %v", tt.expectedResult[i], user)
				}
			}

		})
	}

}
