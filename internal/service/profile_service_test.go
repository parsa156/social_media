package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"social_media/internal/domain"
	"social_media/internal/mocks"
)

// Test 1: Get profile for a nonexistent user.
func TestGetProfileNonExistentUser(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	profileService := NewProfileService(userRepoMock)

	userRepoMock.On("FindByID", "nonexistent").Return(nil, nil)

	profile, err := profileService.GetProfile("nonexistent")
	assert.Nil(t, profile)
	assert.EqualError(t, err, "user not found")
	userRepoMock.AssertExpectations(t)
}

// Test 2: Get profile successfully.
func TestGetProfileSuccess(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	profileService := NewProfileService(userRepoMock)

	expectedUser := &domain.User{ID: "user1", Name: "Alice"}
	userRepoMock.On("FindByID", "user1").Return(expectedUser, nil)

	profile, err := profileService.GetProfile("user1")
	assert.NotNil(t, profile)
	assert.Nil(t, err)
	assert.Equal(t, expectedUser, profile)
	userRepoMock.AssertExpectations(t)
}

// Test 3: Update profile with username conflict.
func TestUpdateProfileUsernameConflict(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	profileService := NewProfileService(userRepoMock)

	currentUser := &domain.User{ID: "user1", Name: "Alice", Username: nil}
	conflictingUser := &domain.User{ID: "user2", Name: "Bob", Username: ptr("@bob")}
	userRepoMock.On("FindByID", "user1").Return(currentUser, nil)
	userRepoMock.On("FindByUsername", "@bob").Return(conflictingUser, nil)

	updatedUser, err := profileService.UpdateProfile("user1", "Alice Updated", "bob", "")
	assert.Nil(t, updatedUser)
	assert.EqualError(t, err, "username already used")
	userRepoMock.AssertExpectations(t)
}

// Test 4: Successful profile update.
func TestUpdateProfileSuccess(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	profileService := NewProfileService(userRepoMock)

	currentUser := &domain.User{ID: "user1", Name: "Alice", Username: nil, Password: "oldhash"}
	userRepoMock.On("FindByID", "user1").Return(currentUser, nil)
	userRepoMock.On("FindByUsername", "@aliceNew").Return(nil, nil)
	userRepoMock.On("Update", mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(0).(*domain.User)
		u.Password = "updatedhash"
	})

	updatedUser, err := profileService.UpdateProfile("user1", "Alice New", "aliceNew", "newpassword")
	assert.NotNil(t, updatedUser)
	assert.Nil(t, err)
	assert.Equal(t, "Alice New", updatedUser.Name)
	assert.Equal(t, "@aliceNew", *updatedUser.Username)
	userRepoMock.AssertExpectations(t)
}

// Test 5: Delete profile for a nonexistent user.
func TestDeleteProfileNonExistentUser(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	profileService := NewProfileService(userRepoMock)

	userRepoMock.On("FindByID", "userNonExistent").Return(nil, nil)

	err := profileService.DeleteProfile("userNonExistent")
	assert.EqualError(t, err, "user not found")
	userRepoMock.AssertExpectations(t)
}

// Test 6: Successful profile deletion.
func TestDeleteProfileSuccess(t *testing.T) {
	userRepoMock := new(mocks.UserRepositoryMock)
	profileService := NewProfileService(userRepoMock)

	existingUser := &domain.User{ID: "user1", Name: "Alice"}
	userRepoMock.On("FindByID", "user1").Return(existingUser, nil)
	userRepoMock.On("Delete", existingUser).Return(nil)

	err := profileService.DeleteProfile("user1")
	assert.Nil(t, err)
	userRepoMock.AssertExpectations(t)
}

// Helper to get a pointer to a string.
func ptr(s string) *string {
	return &s
}
