package repository

import (
	"errors"

	"gorm.io/gorm"
	"social_media/internal/domain"
)

type conversationRepository struct {
	db *gorm.DB
}
//NewConversationRepository function
func NewConversationRepository(db *gorm.DB) domain.ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) Create(convo *domain.Conversation) error {
	return r.db.Create(convo).Error
}

func (r *conversationRepository) FindByParticipants(p1, p2 string) (*domain.Conversation, error) {
	// Order participants lexicographically to ensure consistency.
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	var convo domain.Conversation
	if err := r.db.Where("participant1 = ? AND participant2 = ?", p1, p2).First(&convo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &convo, nil
}

func (r *conversationRepository) FindByUser(userID string) ([]*domain.Conversation, error) {
	var convos []*domain.Conversation
	if err := r.db.Where("participant1 = ? OR participant2 = ?", userID, userID).Find(&convos).Error; err != nil {
		return nil, err
	}
	return convos, nil
}

func (r *conversationRepository) FindByID(id string) (*domain.Conversation, error) {
	var convo domain.Conversation
	if err := r.db.Where("id = ?", id).First(&convo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &convo, nil
}
