package mocks

import (
	"github.com/stretchr/testify/mock"
	"social_media/internal/domain"
)

type MessageRepositoryMock struct {
	mock.Mock
}

func (m *MessageRepositoryMock) Create(message *domain.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MessageRepositoryMock) Update(message *domain.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MessageRepositoryMock) Delete(message *domain.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MessageRepositoryMock) FindByConversation(convoID string) ([]*domain.Message, error) {
	args := m.Called(convoID)
	if messages := args.Get(0); messages != nil {
		return messages.([]*domain.Message), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MessageRepositoryMock) FindByID(id string) (*domain.Message, error) {
	args := m.Called(id)
	if mobj := args.Get(0); mobj != nil {
		return mobj.(*domain.Message), args.Error(1)
	}
	return nil, args.Error(1)
}
