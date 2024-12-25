package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"social_media/internal/domain"
)

// ConversationService defines the interface for conversation and messaging operations.
type ConversationService interface {
	// SendMessage creates a conversation (if needed) and sends a message.
	// The recipientIdentifier can be a phone number or a username (with '@').
	SendMessage(senderID, recipientIdentifier, content string) (*domain.Message, error)
	GetConversations(userID string) ([]*domain.Conversation, error)
	GetMessages(convoID string) ([]*domain.Message, error)
	UpdateMessage(senderID, messageID, content string) (*domain.Message, error)
	DeleteMessage(senderID, messageID string) error
}

type conversationService struct {
	convoRepo   domain.ConversationRepository
	messageRepo domain.MessageRepository
	userRepo    domain.UserRepository // Used to lookup recipient details.
}

// NewConversationService creates a new instance of ConversationService.
func NewConversationService(
	convoRepo domain.ConversationRepository,
	messageRepo domain.MessageRepository,
	userRepo domain.UserRepository,
) ConversationService {
	return &conversationService{
		convoRepo:   convoRepo,
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

// SendMessage looks up the recipient by phone or username and sends the message.
func (s *conversationService) SendMessage(senderID, recipientIdentifier, content string) (*domain.Message, error) {
	// Lookup the recipient using the identifier.
	var recipient *domain.User
	var err error
	if recipientIdentifier != "" && recipientIdentifier[0] == '@' {
		// Treat the identifier as a username.
		recipient, err = s.userRepo.FindByUsername(recipientIdentifier)
	} else {
		// Otherwise, assume it's a phone number.
		recipient, err = s.userRepo.FindByPhone(recipientIdentifier)
	}
	if err != nil {
		return nil, err
	}
	if recipient == nil {
		return nil, errors.New("recipient not found")
	}
	recipientID := recipient.ID

	// Order senderID and recipientID lexicographically so that a unique conversation exists per pair.
	p1, p2 := senderID, recipientID
	if p1 > p2 {
		p1, p2 = p2, p1
	}

	// Find or create the conversation.
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

	// Create the message.
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

// GetConversations returns all conversations for the specified user.
func (s *conversationService) GetConversations(userID string) ([]*domain.Conversation, error) {
	return s.convoRepo.FindByUser(userID)
}

// GetMessages returns all messages within the specified conversation.
func (s *conversationService) GetMessages(convoID string) ([]*domain.Message, error) {
	return s.messageRepo.FindByConversation(convoID)
}

// UpdateMessage allows the sender to update their message.
func (s *conversationService) UpdateMessage(senderID, messageID, content string) (*domain.Message, error) {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil || message == nil {
		return nil, errors.New("message not found")
	}
	// Only the original sender can update the message.
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

// DeleteMessage allows the sender to delete their message.
func (s *conversationService) DeleteMessage(senderID, messageID string) error {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil || message == nil {
		return errors.New("message not found")
	}
	// Only the sender can delete the message.
	if message.SenderID != senderID {
		return errors.New("not authorized to delete this message")
	}
	return s.messageRepo.Delete(message)
}
