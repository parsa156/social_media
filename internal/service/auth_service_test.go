package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"social_media/internal/domain"
	"social_media/internal/mocks"
	"social_media/pkg/jwt"
)

// Test 1: Register with a short password.
func TestRegisterWithShortPassword(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	jwtManager := jwt.NewJWTManager("secret", time.Hour*24)
	authService := NewAuthService(userRepoMock, jwtManager)

	user, err := authService.Register("Test User", "1234567890", "", "short")

	assert.Nil(t, user)
	assert.EqualError(t, err, "password must be at least 8 characters")
}

// Test 2: Register duplicate phone.
func TestRegisterDuplicatePhone(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	jwtManager := jwt.NewJWTManager("secret", time.Hour*24)
	authService := NewAuthService(userRepoMock, jwtManager)

	existingUser := &domain.User{ID: "existing-id"}
	userRepoMock.On("FindByPhone", "1234567890").Return(existingUser, nil)

	user, err := authService.Register("Test User", "1234567890", "", "validpassword")

	assert.Nil(t, user)
	assert.EqualError(t, err, "phone already registered")
	userRepoMock.AssertExpectations(t)
}

// Test 3: Register duplicate username.
func TestRegisterDuplicateUsername(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	jwtManager := jwt.NewJWTManager("secret", time.Hour*24)
	authService := NewAuthService(userRepoMock, jwtManager)

	// Phone is new.
	userRepoMock.On("FindByPhone", "0987654321").Return(nil, nil)
	// Username already used.
	userRepoMock.On("FindByUsername", "@existing").Return(&domain.User{ID: "existing-username-id"}, nil)

	user, err := authService.Register("Test User", "0987654321", "existing", "validpassword")

	assert.Nil(t, user)
	assert.EqualError(t, err, "username already used")
	userRepoMock.AssertExpectations(t)
}

// Test 4: Register without username.
func TestRegisterWithoutUsername(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	jwtManager := jwt.NewJWTManager("secret", time.Hour*24)
	authService := NewAuthService(userRepoMock, jwtManager)

	userRepoMock.On("FindByPhone", "1112223333").Return(nil, nil)
	userRepoMock.On("Create", mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(0).(*domain.User)
		u.ID = "new-user-id"
	})

	user, err := authService.Register("Test User", "1112223333", "", "validpassword")

	assert.NotNil(t, user)
	assert.Nil(t, err)
	// When no username is provided, user.Username remains nil.
	assert.Nil(t, user.Username)
	userRepoMock.AssertExpectations(t)
}

// Test 5: Register with username missing '@'.
func TestRegisterWithUsernamePrefix(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	jwtManager := jwt.NewJWTManager("secret", time.Hour*24)
	authService := NewAuthService(userRepoMock, jwtManager)

	userRepoMock.On("FindByPhone", "2223334444").Return(nil, nil)
	userRepoMock.On("FindByUsername", "@newuser").Return(nil, nil)
	userRepoMock.On("Create", mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(0).(*domain.User)
		u.ID = "new-user-id"
	})

	user, err := authService.Register("Test User", "2223334444", "newuser", "validpassword")

	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.Equal(t, "@newuser", *user.Username)
	userRepoMock.AssertExpectations(t)
}

// Test 6: Login with nonexistent user.
func TestLoginWithNonexistentUser(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	jwtManager := jwt.NewJWTManager("secret", time.Hour*24)
	authService := NewAuthService(userRepoMock, jwtManager)

	userRepoMock.On("FindByPhone", "3334445555").Return(nil, nil)

	token, err := authService.Login("3334445555", "anyPassword")
	assert.Empty(t, token)
	assert.EqualError(t, err, "invalid credentials")
	userRepoMock.AssertExpectations(t)
}

// Test 7: Login with incorrect password.
func TestLoginWithIncorrectPassword(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	jwtManager := jwt.NewJWTManager("secret", time.Hour*24)
	authService := NewAuthService(userRepoMock, jwtManager)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:       "user-id",
		Phone:    "4445556666",
		Password: string(hashedPassword),
	}
	userRepoMock.On("FindByPhone", "4445556666").Return(user, nil)

	token, err := authService.Login("4445556666", "wrongpassword")
	assert.Empty(t, token)
	assert.EqualError(t, err, "invalid credentials")
	userRepoMock.AssertExpectations(t)
}

// Test 8: Successful login.
func TestLoginSuccess(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	jwtManager := jwt.NewJWTManager("secret", time.Hour*24)
	authService := NewAuthService(userRepoMock, jwtManager)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("validpassword"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:       "user-id",
		Phone:    "5556667777",
		Password: string(hashedPassword),
		Username: nil,
	}
	userRepoMock.On("FindByPhone", "5556667777").Return(user, nil)

	token, err := authService.Login("5556667777", "validpassword")
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
	userRepoMock.AssertExpectations(t)
}
