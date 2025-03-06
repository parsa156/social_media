package mocks

import (
	"github.com/stretchr/testify/mock"
	"social_media/internal/domain"
)

type ConversationRepositoryMock struct {
	mock.Mock
}

func (m *ConversationRepositoryMock) Create(convo *domain.Conversation) error {
	args := m.Called(convo)
	return args.Error(0)
}

func (m *ConversationRepositoryMock) FindByParticipants(p1, p2 string) (*domain.Conversation, error) {
	args := m.Called(p1, p2)
	if c := args.Get(0); c != nil {
		return c.(*domain.Conversation), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ConversationRepositoryMock) FindByUser(userID string) ([]*domain.Conversation, error) {
	args := m.Called(userID)
	if convos := args.Get(0); convos != nil {
		return convos.([]*domain.Conversation), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ConversationRepositoryMock) FindByID(id string) (*domain.Conversation, error) {
	args := m.Called(id)
	if c := args.Get(0); c != nil {
		return c.(*domain.Conversation), args.Error(1)
	}
	return nil, args.Error(1)
}
