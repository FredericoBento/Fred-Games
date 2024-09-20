package database

import "github.com/FredericoBento/HandGame/internal/models"

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetAll() ([]models.User, error)
	Create(*models.User) error
}
