package database

import "github.com/FredericoBento/HandGame/internal/models"

type UserRepository interface {
	GetAll() ([]models.User, error)
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
}
