package services

import (
	"errors"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/database"
	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/models"
)

var (
	ErrCouldNotGetAllUsers  = errors.New("could not retrive all users from users repository")
	ErrNoUsersFound         = errors.New("no users were found")
	ErrCouldNotCreateLogger = errors.New("could not create user_service logger")
	ErrCouldNotCreateUser   = errors.New("could not create user")
	ErrUserAlreadyExists    = errors.New("user already exists")
)

type UserService struct {
	repo database.UserRepository
	log  *slog.Logger
}

func NewUserService(repo database.UserRepository) *UserService {
	lo, err := logger.NewServiceLogger("UserService", "", false)
	if err != nil {
		slog.Error(ErrCouldNotCreateLogger.Error() + " " + err.Error())
		lo = slog.Default()
	}
	return &UserService{
		repo: repo,
		log:  lo,
	}
}

func (us *UserService) GetAllUsers() ([]models.User, error) {
	users, err := us.repo.GetAll()
	if err != nil {
		us.log.Error(err.Error())
		return nil, ErrCouldNotGetAllUsers
	}

	// if len(users) == 0 {
	// 	return nil, ErrNoUsersFound
	// }

	return users, nil

}

func (us *UserService) UserExists(username string) (bool, error) {
	_, err := us.repo.GetByUsername(username)
	if err != nil {
		us.log.Error(err.Error())
		return false, err
	}
	return true, nil
}

func (us *UserService) CreateUser(user *models.User) error {
	exist, err := us.UserExists(user.Username)
	if err != nil {
		us.log.Error(err.Error())
		return err
	}

	if exist {
		return ErrUserAlreadyExists
	}

	err = us.repo.Create(user)
	if err != nil {
		us.log.Error(err.Error())
		return ErrCouldNotCreateUser
	}
	return nil
}
