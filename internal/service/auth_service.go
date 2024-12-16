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

// AuthService interface
type AuthService interface {
	Register(name, phone, username, password string) (*domain.User, error)
	Login(phone, password string) (string, error)
}

type authService struct {
	userRepo   domain.UserRepository
	jwtManager *jwt.JWTManager
}

// NewAuthService function
func NewAuthService(userRepo domain.UserRepository, jwtManager *jwt.JWTManager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register creates a new user with validation.
func (s *authService) Register(name, phone, username, password string) (*domain.User, error) {
	// Basic validation
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Check if phone is already registered.
	if existing, _ := s.userRepo.FindByPhone(phone); existing != nil {
		return nil, errors.New("phone already registered")
	}

	// If username is provided, check for uniqueness.
	if username != "" {
		if !strings.HasPrefix(username, "@") {
			username = "@" + username
		}
		if existing, _ := s.userRepo.FindByUsername(username); existing != nil {
			return nil, errors.New("username already used")
		}
	} else {
		// if not provided, leave it nil.
		username = ""
	}

	// Generate a UUID for the user.
	userID := uuid.New().String()

	// Hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:        userID,
		Name:      name,
		Phone:     phone,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Set username only if provided.
	if username != "" {
		user.Username = &username
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login validates credentials and returns a JWT token.
func (s *authService) Login(phone, password string) (string, error) {
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return "", err
	}
	return token, nil
}
