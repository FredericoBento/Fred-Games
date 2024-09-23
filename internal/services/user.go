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
)

var (
	ErrCouldNotGetAllUsers  = errors.New("could not retrive all users from users repository")
	ErrNoUsersFound         = errors.New("no users were found")
	ErrCouldNotCreateLogger = errors.New("could not create user_service logger")
	ErrCouldNotCreateUser   = errors.New("could not create user")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrCouldNotGetUser      = errors.New("could not retrieve user from repository")
	ErrUserNotCached        = errors.New("user was not found in cache")
)

type UserService struct {
	repo  repository.UserRepository
	cache sync.Map
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
		repo: repo,
		ttl:  ttl,
		log:  lo,
	}
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
		return false, err
	}

	return true, nil
}

func (us *UserService) CreateUser(ctx context.Context, user *models.User) error {
	exist, err := us.UserExists(ctx, user.Username)
	if err != nil {
		return err
	}

	if exist {
		return ErrUserAlreadyExists
	}

	err = us.repo.Create(ctx, user)
	if err != nil {
		us.log.Error(err.Error())
		return ErrCouldNotCreateUser
	}
	return nil
}

func (us *UserService) ComparePassword(password1 string, password2 string) (bool, error) {
	// Hash if needed
	return password1 == password2, nil

}

func (us *UserService) evictAfter(key string, ttl time.Duration) {
	time.Sleep(ttl)
	us.cache.Delete(key)
}

func (us *UserService) getUserInCache(username string) (*models.User, error) {
	if cachedUser, ok := us.cache.Load(username); ok {
		if user, valid := cachedUser.(*models.User); valid {
			return user, nil
		}
	}
	return nil, ErrUserNotCached
}
