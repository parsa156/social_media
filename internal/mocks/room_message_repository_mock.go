package mocks

import (
	"github.com/stretchr/testify/mock"
	"social_media/internal/domain"
)

type RoomMessageRepositoryMock struct {
	mock.Mock
}

func (m *RoomMessageRepositoryMock) Create(message *domain.RoomMessage) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *RoomMessageRepositoryMock) Update(message *domain.RoomMessage) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *RoomMessageRepositoryMock) Delete(messageID string) error {
	args := m.Called(messageID)
	return args.Error(0)
}

func (m *RoomMessageRepositoryMock) FindByRoom(roomID string) ([]*domain.RoomMessage, error) {
	args := m.Called(roomID)
	if messages := args.Get(0); messages != nil {
		return messages.([]*domain.RoomMessage), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *RoomMessageRepositoryMock) FindByID(messageID string) (*domain.RoomMessage, error) {
	args := m.Called(messageID)
	if mobj := args.Get(0); mobj != nil {
		return mobj.(*domain.RoomMessage), args.Error(1)
	}
	return nil, args.Error(1)
}
