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

// Test 4: Delete room unauthorized.
func TestDeleteRoomUnauthorized(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	room := &domain.Room{ID: "room1", Name: "Test Room", OwnerID: "owner1"}
	roomRepoMock.On("FindByID", "room1").Return(room, nil)
	membershipRepoMock.On("GetMemberRole", "room1", "user2").Return(domain.RoleMember, nil)

	err := roomService.DeleteRoom("room1", "user2")
	assert.EqualError(t, err, "not authorized to delete room")
	roomRepoMock.AssertExpectations(t)
	membershipRepoMock.AssertExpectations(t)
}

// Test 5: Add member already exists.
func TestAddMemberAlreadyExists(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	room := &domain.Room{ID: "room1", Type: domain.RoomTypeGroup}
	roomRepoMock.On("FindByID", "room1").Return(room, nil)
	membershipRepoMock.On("GetMemberRole", "room1", "userX").Return(domain.RoleMember, nil)

	err := roomService.AddMember("room1", "owner1", "userX")
	assert.EqualError(t, err, "user already a member")
	roomRepoMock.AssertExpectations(t)
	membershipRepoMock.AssertExpectations(t)
}

// Test 6: Remove member unauthorized.
func TestRemoveMemberUnauthorized(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	// Only set expectation for the requester (user3) since the code checks that role.
	membershipRepoMock.On("GetMemberRole", "room1", "user3").Return(domain.RoleMember, nil)

	err := roomService.RemoveMember("room1", "user3", "user2")
	assert.EqualError(t, err, "not authorized to remove member")
	membershipRepoMock.AssertExpectations(t)
}


// Test 7: Promote member unauthorized.
func TestPromoteMemberUnauthorized(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	// Requester is not owner.
	membershipRepoMock.On("GetMemberRole", "room1", "user3").Return(domain.RoleMember, nil)

	err := roomService.PromoteMember("room1", "user3", "user2")
	assert.EqualError(t, err, "only owner can promote member")
	membershipRepoMock.AssertExpectations(t)
}

// Test 8: Ban member unauthorized.
func TestBanMemberUnauthorized(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	// Requester is not owner/admin.
	membershipRepoMock.On("GetMemberRole", "room1", "user3").Return(domain.RoleMember, nil)

	err := roomService.BanMember("room1", "user3", "user2")
	assert.EqualError(t, err, "not authorized to ban member")
	membershipRepoMock.AssertExpectations(t)
}

// Test 9: Send room message by a banned user.
func TestSendRoomMessageBannedUser(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	membershipRepoMock.On("IsUserBanned", "room1", "user1").Return(true, nil)

	msg, err := roomService.SendMessage("room1", "user1", "Hello in room")
	assert.Nil(t, msg)
	assert.EqualError(t, err, "you are banned from this room")
	membershipRepoMock.AssertExpectations(t)
}
//Test 10 Send Room Message Success
func TestSendRoomMessageSuccess(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	// Arrange: User is not banned, and room exists.
	membershipRepoMock.On("IsUserBanned", "room1", "user1").Return(false, nil)
	room := &domain.Room{ID: "room1", Type: domain.RoomTypeGroup}
	roomRepoMock.On("FindByID", "room1").Return(room, nil)
	roomMessageRepoMock.On("Create", mock.AnythingOfType("*domain.RoomMessage")).Return(nil).Run(func(args mock.Arguments) {
		msg := args.Get(0).(*domain.RoomMessage)
		msg.ID = "msg1"
	})

	// Act: User sends a message.
	msg, err := roomService.SendMessage("room1", "user1", "Hello Room!")

	// Assert: The message is created successfully.
	assert.NotNil(t, msg)
	assert.Equal(t, "msg1", msg.ID)
	assert.Nil(t, err)
	membershipRepoMock.AssertExpectations(t)
	roomRepoMock.AssertExpectations(t)
	roomMessageRepoMock.AssertExpectations(t)
}
//Test 11 Unban Member Success
func TestUnbanMemberSuccess(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	// Arrange: Requester is admin.
	membershipRepoMock.On("GetMemberRole", "room1", "admin1").Return(domain.RoleAdmin, nil)
	membershipRepoMock.On("UpdateMemberRole", "room1", "user2", domain.RoleMember).Return(nil)

	// Act: Admin unbans a member.
	err := roomService.UnbanMember("room1", "admin1", "user2")

	// Assert: No error is returned.
	assert.Nil(t, err)
	membershipRepoMock.AssertExpectations(t)
}

//Test 12 Delete Room Success
func TestDeleteRoomSuccess(t *testing.T) {
	roomRepoMock := new(mocks.RoomRepositoryMock)
	membershipRepoMock := new(mocks.RoomMembershipRepositoryMock)
	roomMessageRepoMock := new(mocks.RoomMessageRepositoryMock)
	roomService := NewRoomService(roomRepoMock, membershipRepoMock, roomMessageRepoMock)

	// Arrange: Room exists and requester is owner.
	room := &domain.Room{ID: "room1", OwnerID: "owner1"}
	roomRepoMock.On("FindByID", "room1").Return(room, nil)
	membershipRepoMock.On("GetMemberRole", "room1", "owner1").Return(domain.RoleOwner, nil)
	roomRepoMock.On("Delete", "room1").Return(nil)

	// Act: Owner deletes the room.
	err := roomService.DeleteRoom("room1", "owner1")

	// Assert: No error occurs.
	assert.Nil(t, err)
	roomRepoMock.AssertExpectations(t)
	membershipRepoMock.AssertExpectations(t)
}

