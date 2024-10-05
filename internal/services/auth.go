package services

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/models"
	"github.com/google/uuid"
)

var (
	ErrCouldNotComparePassword = errors.New("could not compare password")
	ErrCouldNotFindUser        = errors.New("could not find user")
	ErrIncorrectCredentials    = errors.New("wrong credentials provided")
	ErrSessionExpired          = errors.New("session has expired")
	ErrTokenDoesNotExists      = errors.New("token invalid, does not exist")
	ErrNoToken                 = errors.New("token not found")
)

const (
	sessionExpiryTime = 120 * time.Minute
	cookieName        = "session_token"
)

type Session struct {
	username string
	expiry   time.Time
}

type AuthService struct {
	sessions    map[string]Session
	userService *UserService
	mu          sync.Mutex
	log         *slog.Logger
}

func NewAuthService(userService *UserService) *AuthService {
	if userService == nil {
		log.Fatal("no user service provided")
	}

	lo, err := logger.NewServiceLogger("AuthService", "", true)
	if err != nil {
		lo = slog.Default()
	}

	return &AuthService{
		sessions:    make(map[string]Session),
		userService: userService,
		mu:          sync.Mutex{},
		log:         lo,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, username string, password string) (*models.User, error) {
	user, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		s.log.Error(err.Error())
		return nil, ErrCouldNotFindUser
	}

	passwordCheck := s.userService.ComparePassword(user.Password, password)
	// if err != nil {
	// 	s.log.Error(err.Error())
	// 	return nil, ErrCouldNotComparePassword
	// }

	if passwordCheck != true {
		return nil, ErrIncorrectCredentials
	}

	return user, nil

}

func (s *AuthService) CreateSession(user *models.User) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionToken := uuid.NewString()
	session := Session{
		username: user.Username,
		expiry:   time.Now().Add(sessionExpiryTime),
	}

	s.sessions[sessionToken] = session

	return sessionToken, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[token]
	if !exists {
		return nil, ErrTokenDoesNotExists
	}

	if session.IsExpired() {
		return nil, ErrSessionExpired
	}

	user, err := s.userService.GetUserByUsername(ctx, session.username)
	if err != nil {
		s.log.Error(err.Error())
		return nil, ErrCouldNotFindUser
	}

	return user, nil
}

func (s *AuthService) DestroySession(ctx context.Context, token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, token)
}

func (s *AuthService) GetCookieName() string {
	return cookieName
}

func (s *Session) IsExpired() bool {
	return s.expiry.Before(time.Now())
}

func (s *AuthService) IsAdmin(username string) bool {
	return username == "fred"
}

func (s *AuthService) GetToken(r *http.Request) (string, error) {
	c, err := r.Cookie(s.GetCookieName())

	if err == nil && c != nil {
		return c.Value, nil
	}

	token := r.Header.Get("Authorization")

	if token == "" {
		return "", ErrNoToken
	}

	return token, nil

}

func (s *AuthService) IsLogged(r *http.Request) bool {
	token, err := s.GetToken(r)
	if err != nil {
		s.log.Error(err.Error())
		return false
	}
	_, err = s.ValidateSession(context.Background(), token)
	if err != nil {
		return false
	}
	return true
}
