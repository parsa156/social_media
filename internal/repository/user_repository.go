package repository

import (
	"errors"

	"gorm.io/gorm"
	"social_media/internal/domain"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByPhone(phone string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
