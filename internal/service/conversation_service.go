package service

import (
	"errors"
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

//NewConversationService : function
func NewConversationService(convoRepo domain.ConversationRepository, messageRepo domain.MessageRepository) ConversationService {
	return &conversationService{
		convoRepo:   convoRepo,
		messageRepo: messageRepo,
	}
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
	message := &domain.Message{
		ID:             uuid.New().String(),
		ConversationID: convo.ID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}
	return message, nil
}

func (s *conversationService) GetConversations(userID string) ([]*domain.Conversation, error) {
	return s.convoRepo.FindByUser(userID)
}

func (s *conversationService) GetMessages(convoID string) ([]*domain.Message, error) {
	return s.messageRepo.FindByConversation(convoID)
}

func (s *conversationService) UpdateMessage(senderID, messageID, content string) (*domain.Message, error) {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil || message == nil {
		return nil, errors.New("message not found")
	}
	// Only the sender can update their message.
	if message.SenderID != senderID {
		return nil, errors.New("not authorized to update this message")
	}
	message.Content = content
	message.UpdatedAt = time.Now()
	if err := s.messageRepo.Update(message); err != nil {
		return nil, err
	}
	return message, nil
}

func (s *conversationService) DeleteMessage(senderID, messageID string) error {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil || message == nil {
		return errors.New("message not found")
	}
	// Only the sender can delete their message.
	if message.SenderID != senderID {
		return errors.New("not authorized to delete this message")
	}
	return s.messageRepo.Delete(message)
}
