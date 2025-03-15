package mocks

import (
	"github.com/stretchr/testify/mock"
	"social_media/internal/domain"
)

type RoomRepositoryMock struct {
	mock.Mock
}

func (m *RoomRepositoryMock) Create(room *domain.Room) error {
	args := m.Called(room)
	return args.Error(0)
}

func (m *RoomRepositoryMock) Update(room *domain.Room) error {
	args := m.Called(room)
	return args.Error(0)
}

func (m *RoomRepositoryMock) Delete(roomID string) error {
	args := m.Called(roomID)
	return args.Error(0)
}

func (m *RoomRepositoryMock) FindByID(roomID string) (*domain.Room, error) {
	args := m.Called(roomID)
	if r := args.Get(0); r != nil {
		return r.(*domain.Room), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *RoomRepositoryMock) FindByUsername(username string) (*domain.Room, error) {
	args := m.Called(username)
	if r := args.Get(0); r != nil {
		return r.(*domain.Room), args.Error(1)
	}
	return nil, args.Error(1)
}
