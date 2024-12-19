package domain

import "time"

// Message represents an individual message in a conversation.
type Message struct {
	ID             string    `gorm:"type:uuid;primaryKey" json:"id"`
	ConversationID string    `gorm:"type:uuid;not null" json:"conversation_id"`
	SenderID       string    `gorm:"type:uuid;not null" json:"sender_id"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}


// MessageRepository defines the methods for message persistence.
type MessageRepository interface {
	Create(message *Message) error
	Update(message *Message) error
	Delete(message *Message) error
	FindByConversation(convoID string) ([]*Message, error)
	FindByID(id string) (*Message, error)
}