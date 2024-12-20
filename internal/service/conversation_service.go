package service

import (
	"time"

	"github.com/google/uuid"
	"social_media/internal/domain"
)

//ConversationService : interface
type ConversationService interface {
	// SendMessage creates a conversation if needed and sends a message.
	SendMessage(senderID, recipientID, content string) (*domain.Message, error)
	GetConversations(userID string) ([]*domain.Conversation, error)
	GetMessages(convoID string) ([]*domain.Message, error)
	UpdateMessage(senderID, messageID, content string) (*domain.Message, error)
	DeleteMessage(senderID, messageID string) error
}

type conversationService struct {
	convoRepo   domain.ConversationRepository
	messageRepo domain.MessageRepository
}


func (s *conversationService) SendMessage(senderID, recipientID, content string) (*domain.Message, error) {
	// Order participant IDs so the conversation is unique.
	p1, p2 := senderID, recipientID
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	// Check if conversation exists.
	convo, err := s.convoRepo.FindByParticipants(p1, p2)
	if err != nil {
		return nil, err
	}
	if convo == nil {
		convo = &domain.Conversation{
			ID:           uuid.New().String(),
			Participant1: p1,
			Participant2: p2,
			CreatedAt:    time.Now(),
		}
		if err := s.convoRepo.Create(convo); err != nil {
			return nil, err
		}
	}
}	