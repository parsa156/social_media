package repository

import (
	"errors"

	"gorm.io/gorm"
	"social_media/internal/domain"
)

type messageRepository struct {
	db *gorm.DB
}
//NewMessageRepository function
func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *domain.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) Update(message *domain.Message) error {
	return r.db.Save(message).Error
}

func (r *messageRepository) Delete(message *domain.Message) error {
	return r.db.Delete(message).Error
}

func (r *messageRepository) FindByConversation(convoID string) ([]*domain.Message, error) {
	var messages []*domain.Message
	if err := r.db.Where("conversation_id = ?", convoID).Order("created_at asc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) FindByID(id string) (*domain.Message, error) {
	var message domain.Message
	if err := r.db.Where("id = ?", id).First(&message).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}
