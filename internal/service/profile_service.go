package service

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"social_media/internal/domain"
)
// ProfileService interface 
type ProfileService interface {
	GetProfile(userID string) (*domain.User, error)
	UpdateProfile(userID, name, username, password string) (*domain.User, error)
	DeleteProfile(userID string) error
}

type profileService struct {
	userRepo domain.UserRepository
}
// NewProfileService function
func NewProfileService(userRepo domain.UserRepository) ProfileService {
	return &profileService{
		userRepo: userRepo,
	}
}

// GetProfile returns the user profile by ID.
func (s *profileService) GetProfile(userID string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// UpdateProfile allows updating the name, username, and password.
func (s *profileService) UpdateProfile(userID, name, username, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Update name if provided.
	if name != "" {
		user.Name = name
	}

	// Update username if provided.
	if username != "" {
		// Prepend '@' if missing.
		if !strings.HasPrefix(username, "@") {
			username = "@" + username
		}
		// Check for uniqueness.
		existing, _ := s.userRepo.FindByUsername(username)
		// Allow update if the found user is not the current one.
		if existing != nil && existing.ID != user.ID {
			return nil, errors.New("username already used")
		}
		user.Username = &username
	}

	// Update password if provided.
	if password != "" {
		if len(password) < 8 {
			return nil, errors.New("password must be at least 8 characters")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteProfile deletes the user account.
func (s *profileService) DeleteProfile(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	return s.userRepo.Delete(user)
}
