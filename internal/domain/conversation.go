package domain

import "time"

// Conversation represents a conversation between two users.
type Conversation struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	Participant1 string    `gorm:"type:uuid;not null" json:"participant1"`
	Participant2 string    `gorm:"type:uuid;not null" json:"participant2"`
	CreatedAt    time.Time `json:"created_at"`
}



// ConversationRepository defines the methods for conversation persistence.
type ConversationRepository interface {
	Create(convo *Conversation) error
	FindByParticipants(p1, p2 string) (*Conversation, error)
	FindByUser(userID string) ([]*Conversation, error)
	FindByID(id string) (*Conversation, error)
}

