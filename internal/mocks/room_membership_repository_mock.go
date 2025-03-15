package mocks

import (
	"github.com/stretchr/testify/mock"
	"social_media/internal/domain"
)

type RoomMembershipRepositoryMock struct {
	mock.Mock
}

func (m *RoomMembershipRepositoryMock) AddMember(membership *domain.RoomMembership) error {
	args := m.Called(membership)
	return args.Error(0)
}

func (m *RoomMembershipRepositoryMock) UpdateMemberRole(roomID, userID string, role domain.RoomMembershipRole) error {
	args := m.Called(roomID, userID, role)
	return args.Error(0)
}

func (m *RoomMembershipRepositoryMock) RemoveMember(roomID, userID string) error {
	args := m.Called(roomID, userID)
	return args.Error(0)
}

func (m *RoomMembershipRepositoryMock) GetMembers(roomID string) ([]*domain.RoomMembership, error) {
	args := m.Called(roomID)
	if memberships := args.Get(0); memberships != nil {
		return memberships.([]*domain.RoomMembership), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *RoomMembershipRepositoryMock) IsUserBanned(roomID, userID string) (bool, error) {
	args := m.Called(roomID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *RoomMembershipRepositoryMock) GetMemberRole(roomID, userID string) (domain.RoomMembershipRole, error) {
	args := m.Called(roomID, userID)
	return args.Get(0).(domain.RoomMembershipRole), args.Error(1)
}
