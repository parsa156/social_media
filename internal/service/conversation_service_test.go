package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"social_media/internal/domain"
	"social_media/internal/mocks"
)

// Test 1: Send message with recipient not found.
func TestSendMessageRecipientNotFound(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	convoRepoMock := new(mocks.ConversationRepositoryMock)
	messageRepoMock := new(mocks.MessageRepositoryMock)
	convoService := NewConversationService(convoRepoMock, messageRepoMock, userRepoMock)

	// Simulate recipient not found (using phone)
	userRepoMock.On("FindByPhone", "9998887777").Return(nil, nil)

	msg, err := convoService.SendMessage("sender1", "9998887777", "Hello!")
	assert.Nil(t, msg)
	assert.EqualError(t, err, "recipient not found")
	userRepoMock.AssertExpectations(t)
}

// Test 2: Send message creating a new conversation.
func TestSendMessageNewConversation(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	convoRepoMock := new(mocks.ConversationRepositoryMock)
	messageRepoMock := new(mocks.MessageRepositoryMock)
	convoService := NewConversationService(convoRepoMock, messageRepoMock, userRepoMock)

	// Recipient found by phone.
	recipient := &domain.User{ID: "recipient1"}
	userRepoMock.On("FindByPhone", "1231231234").Return(recipient, nil)

	// Order sender and recipient lexicographically.
	p1, p2 := "sender1", "recipient1"
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	// No conversation exists.
	convoRepoMock.On("FindByParticipants", p1, p2).Return(nil, nil)
	convoRepoMock.On("Create", mock.AnythingOfType("*domain.Conversation")).Return(nil)
	messageRepoMock.On("Create", mock.AnythingOfType("*domain.Message")).Return(nil)

	msg, err := convoService.SendMessage("sender1", "1231231234", "Hi there!")
	assert.NotNil(t, msg)
	assert.Nil(t, err)
	userRepoMock.AssertExpectations(t)
	convoRepoMock.AssertExpectations(t)
	messageRepoMock.AssertExpectations(t)
}

// Test 3: Send message using an existing conversation.
func TestSendMessageExistingConversation(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	convoRepoMock := new(mocks.ConversationRepositoryMock)
	messageRepoMock := new(mocks.MessageRepositoryMock)
	convoService := NewConversationService(convoRepoMock, messageRepoMock, userRepoMock)

	recipient := &domain.User{ID: "recipient1"}
	userRepoMock.On("FindByPhone", "1231231234").Return(recipient, nil)

	p1, p2 := "sender1", "recipient1"
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	existingConvo := &domain.Conversation{ID: "convo1", Participant1: p1, Participant2: p2, CreatedAt: time.Now()}
	convoRepoMock.On("FindByParticipants", p1, p2).Return(existingConvo, nil)
	messageRepoMock.On("Create", mock.AnythingOfType("*domain.Message")).Return(nil)

	msg, err := convoService.SendMessage("sender1", "1231231234", "Hi again!")
	assert.NotNil(t, msg)
	assert.Nil(t, err)
	userRepoMock.AssertExpectations(t)
	convoRepoMock.AssertExpectations(t)
	messageRepoMock.AssertExpectations(t)
}

// Test 4: Update message unauthorized.
func TestUpdateMessageUnauthorized(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	convoRepoMock := new(mocks.ConversationRepositoryMock)
	messageRepoMock := new(mocks.MessageRepositoryMock)
	convoService := NewConversationService(convoRepoMock, messageRepoMock, userRepoMock)

	existingMessage := &domain.Message{ID: "msg1", SenderID: "sender1", Content: "Original", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	messageRepoMock.On("FindByID", "msg1").Return(existingMessage, nil)

	updatedMsg, err := convoService.UpdateMessage("anotherSender", "msg1", "Updated content")
	assert.Nil(t, updatedMsg)
	assert.EqualError(t, err, "not authorized to update this message")
	messageRepoMock.AssertExpectations(t)
}

// Test 5: Successful message update.
func TestUpdateMessageSuccess(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	convoRepoMock := new(mocks.ConversationRepositoryMock)
	messageRepoMock := new(mocks.MessageRepositoryMock)
	convoService := NewConversationService(convoRepoMock, messageRepoMock, userRepoMock)

	existingMessage := &domain.Message{ID: "msg1", SenderID: "sender1", Content: "Original", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	messageRepoMock.On("FindByID", "msg1").Return(existingMessage, nil)
	messageRepoMock.On("Update", mock.AnythingOfType("*domain.Message")).Return(nil).Run(func(args mock.Arguments) {
		m := args.Get(0).(*domain.Message)
		m.UpdatedAt = time.Now()
	})

	updatedMsg, err := convoService.UpdateMessage("sender1", "msg1", "Updated content")
	assert.NotNil(t, updatedMsg)
	assert.Nil(t, err)
	assert.Equal(t, "Updated content", updatedMsg.Content)
	messageRepoMock.AssertExpectations(t)
}

// Test 6: Delete message unauthorized.
func TestDeleteMessageUnauthorized(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	convoRepoMock := new(mocks.ConversationRepositoryMock)
	messageRepoMock := new(mocks.MessageRepositoryMock)
	convoService := NewConversationService(convoRepoMock, messageRepoMock, userRepoMock)

	existingMessage := &domain.Message{ID: "msg1", SenderID: "sender1", Content: "To be deleted", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	messageRepoMock.On("FindByID", "msg1").Return(existingMessage, nil)

	err := convoService.DeleteMessage("anotherSender", "msg1")
	assert.EqualError(t, err, "not authorized to delete this message")
	messageRepoMock.AssertExpectations(t)
}

// Test 7: Successful message deletion.
func TestDeleteMessageSuccess(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	convoRepoMock := new(mocks.ConversationRepositoryMock)
	messageRepoMock := new(mocks.MessageRepositoryMock)
	convoService := NewConversationService(convoRepoMock, messageRepoMock, userRepoMock)

	existingMessage := &domain.Message{ID: "msg1", SenderID: "sender1", Content: "To be deleted", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	messageRepoMock.On("FindByID", "msg1").Return(existingMessage, nil)
	messageRepoMock.On("Delete", existingMessage).Return(nil)

	err := convoService.DeleteMessage("sender1", "msg1")
	assert.Nil(t, err)
	messageRepoMock.AssertExpectations(t)
}
