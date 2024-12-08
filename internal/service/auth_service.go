package service

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

	"social_media/internal/domain"
	"social_media/pkg/jwt"
)

type AuthService interface {
	Register(name, phone, username, password string) (*domain.User, error)
	Login(phone, password string) (string, error)
}

type authService struct {
	userRepo   domain.UserRepository
	jwtManager *jwt.JWTManager
}

// NewAuthService returns a new AuthService.
func NewAuthService(userRepo domain.UserRepository, jwtManager *jwt.JWTManager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register creates a new user.
// - If the username does not start with '@', it is added.
// - Generates a unique UUID for the user.
func (s *authService) Register(name, phone, username, password string) (*domain.User, error) {
	// Check if phone is already registered.
	if existing, _ := s.userRepo.FindByPhone(phone); existing != nil {
		return nil, errors.New("phone already registered")
	}

	// Prepend '@' to the username if needed.
	if !strings.HasPrefix(username, "@") {
		username = "@" + username
	}

	// Generate a unique UUID.
	uuidCode := uuid.New().String()

	// Hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		UUID:      uuidCode,
		Name:      name,
		Phone:     phone,
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login validates the phone and password and returns a JWT token.
func (s *authService) Login(phone, password string) (string, error) {
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// Compare the hashed password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate a JWT token.
	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

