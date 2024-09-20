package mock

import "github.com/FredericoBento/HandGame/internal/models"

type MockUserRepository struct {
	GetByUsernameResult *models.User
	GetByUsernameError  error

	CreateError error

	GetAllResult []models.User
	GetAllError  error
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	if m.GetByUsernameError != nil {
		return nil, m.GetByUsernameError
	}
	return m.GetByUsernameResult, nil
}

func (m *MockUserRepository) Create(user *models.User) error {
	return m.CreateError
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	return m.GetAllResult, m.GetAllError
}
