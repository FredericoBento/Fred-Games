package mock

import "context"

type MockUserService struct {
	UserExistsResult bool
	UserExistsError  error
}

func (m *MockUserService) UserExists(ctx context.Context, username string) (bool, error) {
	return m.UserExistsResult, m.UserExistsError
}
