package repository

import (
	"context"

	"github.com/FredericoBento/HandGame/internal/models"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, user *models.User) error
}
