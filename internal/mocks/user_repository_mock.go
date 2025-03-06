package mocks

import (
	"github.com/stretchr/testify/mock"
	"social_media/internal/domain"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindByPhone(phone string) (*domain.User, error) {
	args := m.Called(phone)
	if u := args.Get(0); u != nil {
		return u.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepositoryMock) FindByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if u := args.Get(0); u != nil {
		return u.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepositoryMock) FindByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if u := args.Get(0); u != nil {
		return u.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepositoryMock) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
