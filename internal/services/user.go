package services

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/FredericoBento/HandGame/internal/database/repository"
	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCouldNotGetAllUsers  = errors.New("could not retrive all users from users repository")
	ErrNoUsersFound         = errors.New("no users were found")
	ErrCouldNotCreateLogger = errors.New("could not create user_service logger")
	ErrCouldNotCreateUser   = errors.New("could not create user")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrCouldNotGetUser      = errors.New("could not retrieve user from repository")
	ErrUserNotCached        = errors.New("user was not found in cache")
	ErrUserExistsFailed     = errors.New("failed to check if user exists")
	ErrInvalidLogger        = errors.New("invalid logger passed")
	ErrCouldNotContactDB    = errors.New("call to repository resulted in a error, could not contact db")
	ErrCouldNotHashPassword = errors.New("could not hash password")
)

type UserService struct {
	name  string
	repo  repository.UserRepository
	cache *sync.Map
	ttl   time.Duration
	log   *slog.Logger
}

func NewUserService(repo repository.UserRepository, ttl time.Duration) *UserService {
	lo, err := logger.NewServiceLogger("UserService", "", false)
	if err != nil {
		slog.Error(ErrCouldNotCreateLogger.Error() + " " + err.Error())
		lo = slog.Default()
	}

	return &UserService{
		repo:  repo,
		cache: &sync.Map{},
		ttl:   ttl,
		log:   lo,
	}
}

// func (*UserService) Start() error {
// }

// func (*UserService) Stop() error {
// }

// func (*UserService) Resume() error {
// }

// func (*UserService) GetName() string {
// }

// func (*UserService) GetRoute() string {
// }

// func (*UserService) GetStatus() Status {
// }

// func (*UserService) GetLogs() ([]logger.PrettyLogs, error) {
// }

func (us *UserService) ChangeLogger(logger *slog.Logger) error {
	if logger == nil {
		return ErrInvalidLogger
	}
	us.log = logger
	return nil
}
func (us *UserService) GetAllUsers() ([]models.User, error) {
	users, err := us.repo.GetAll(context.TODO())
	if err != nil {
		us.log.Error(err.Error())
		return nil, ErrCouldNotGetAllUsers
	}

	return users, nil

}

func (us *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := us.getUserInCache(username)
	if err == nil && user != nil {
		return user, nil
	}

	user, err = us.repo.GetByUsername(ctx, username)
	if err != nil {
		us.log.Error(err.Error())
		return nil, ErrCouldNotGetUser
	}

	us.cache.Store(user.Username, user)
	go us.evictAfter(user.Username, us.ttl)

	return user, nil
}

func (us *UserService) UserExists(ctx context.Context, username string) (bool, error) {
	_, err := us.repo.GetByUsername(ctx, username)
	if err == sql.ErrNoRows {
		// us.log.Error(err.Error())
		return false, nil
	}

	if err != nil {
		us.log.Error(err.Error())
		return false, ErrCouldNotContactDB
	}

	return true, nil
}

func (us *UserService) CreateUser(ctx context.Context, user *models.User) error {
	exist, err := us.UserExists(ctx, user.Username)
	if err != nil {
		us.log.Error(err.Error())
		return ErrUserExistsFailed
	}

	if exist {
		return ErrUserAlreadyExists
	}

	hashed, err := us.HashPassword(user.Password)
	if err != nil {
		return ErrCouldNotCreateUser
	}

	user.Password = hashed

	err = us.repo.Create(ctx, user)
	if err != nil {
		us.log.Error(err.Error())
		return ErrCouldNotCreateUser
	}
	return nil
}

func (us *UserService) ComparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (us *UserService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		us.log.Error(err.Error())
		return "", ErrCouldNotHashPassword
	}

	return string(hashed), nil
}

func (us *UserService) evictAfter(key string, ttl time.Duration) {
	time.Sleep(ttl)
	us.cache.Delete(key)
}

func (us *UserService) getUserInCache(username string) (*models.User, error) {
	if cachedUser, ok := us.cache.Load(username); ok {
		if user, valid := cachedUser.(*models.User); valid {
			us.log.Info("user cache hit")
			return user, nil
		}
	}
	us.log.Info("user cache miss")
	return nil, ErrUserNotCached
}
