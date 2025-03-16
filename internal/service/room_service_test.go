package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"social_media/internal/domain"
	"social_media/internal/mocks"
)

// Test 1: Create room successfully.
func TestCreateRoomSuccess(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	roomRepoMock.On("Create", mock.AnythingOfType("*domain.Room")).Return(nil).Run(func(args mock.Arguments) {
		r := args.Get(0).(*domain.Room)
		r.ID = "room1"
	})
	membershipRepoMock.On("AddMember", mock.AnythingOfType("*domain.RoomMembership")).Return(nil)

	room, err := roomService.CreateRoom("owner1", "Test Room", "testroom", domain.RoomTypeGroup)
	assert.NotNil(t, room)
	assert.Nil(t, err)
	assert.Equal(t, "Test Room", room.Name)
	roomRepoMock.AssertExpectations(t)
	membershipRepoMock.AssertExpectations(t)
}

// Test 2: Update room unauthorized.
func TestUpdateRoomUnauthorized(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	room := &domain.Room{ID: "room1", Name: "Test Room", OwnerID: "owner1", UpdatedAt: time.Now()}
	roomRepoMock.On("FindByID", "room1").Return(room, nil)
	membershipRepoMock.On("GetMemberRole", "room1", "user2").Return(domain.RoleMember, nil)

	updatedRoom, err := roomService.UpdateRoom("room1", "user2", "New Room Name", "newusername")
	assert.Nil(t, updatedRoom)
	assert.EqualError(t, err, "not authorized to update room")
	roomRepoMock.AssertExpectations(t)
	membershipRepoMock.AssertExpectations(t)
}

// Test 3: Successful room update.
func TestUpdateRoomSuccess(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	room := &domain.Room{ID: "room1", Name: "Test Room", OwnerID: "owner1", UpdatedAt: time.Now()}
	roomRepoMock.On("FindByID", "room1").Return(room, nil)
	membershipRepoMock.On("GetMemberRole", "room1", "owner1").Return(domain.RoleOwner, nil)
	roomRepoMock.On("Update", room).Return(nil).Run(func(args mock.Arguments) {
		r := args.Get(0).(*domain.Room)
		r.Name = "Updated Room Name"
	})

	updatedRoom, err := roomService.UpdateRoom("room1", "owner1", "Updated Room Name", "updatedusername")
	assert.NotNil(t, updatedRoom)
	assert.Nil(t, err)
	assert.Equal(t, "Updated Room Name", updatedRoom.Name)
	roomRepoMock.AssertExpectations(t)
	membershipRepoMock.AssertExpectations(t)
}

