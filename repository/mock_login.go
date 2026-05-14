package repository

import (
	"IndulgenceMealPlan/model"

	"github.com/stretchr/testify/mock"
)

type MockLoginRepository struct {
	mock.Mock
}

func (m *MockLoginRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockLoginRepository) GetUserByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockLoginRepository) UpdateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockLoginRepository) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
