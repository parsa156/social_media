package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"social_media/internal/domain"
)

type RoomService interface {
	CreateRoom(ownerID, name, username string, roomType domain.RoomType) (*domain.Room, error)
	UpdateRoom(roomID, updaterID, newName, newUsername string) (*domain.Room, error)
	DeleteRoom(roomID, requesterID string) error
	AddMember(roomID, requesterID, userID string) error
	RemoveMember(roomID, requesterID, userID string) error
	PromoteMember(roomID, requesterID, userID string) error
	BanMember(roomID, requesterID, userID string) error
	UnbanMember(roomID, requesterID, userID string) error
	SendMessage(roomID, senderID, content string) (*domain.RoomMessage, error)
	DeleteMessage(roomID, requesterID, messageID string) error
	GetMessages(roomID string) ([]*domain.RoomMessage, error)
	GetMembers(roomID string) ([]*domain.RoomMembership, error)
}

type roomService struct {
	roomRepo       domain.RoomRepository
	membershipRepo domain.RoomMembershipRepository
	messageRepo    domain.RoomMessageRepository
}

func NewRoomService(
	roomRepo domain.RoomRepository,
	membershipRepo domain.RoomMembershipRepository,
	messageRepo domain.RoomMessageRepository,
)

func (s *roomService) CreateRoom(ownerID, name, username string, roomType domain.RoomType) (*domain.Room, error) {
	room := &domain.Room{
		ID:        uuid.New().String(),
		Name:      name,
		Username:  nil,
		Type:      roomType,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if username != "" {
		room.Username = &username
	}
	if err := s.roomRepo.Create(room); err != nil {
		return nil, err
	}
	// Add owner as a member with owner role.
	membership := &domain.RoomMembership{
		RoomID:    room.ID,
		UserID:    ownerID,
		Role:      domain.RoleOwner,
		CreatedAt: time.Now(),
	}
	if err := s.membershipRepo.AddMember(membership); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *roomService) UpdateRoom(roomID, updaterID, newName, newUsername string) (*domain.Room, error) {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil || room == nil {
		return nil, errors.New("room not found")
	}
	role, err := s.membershipRepo.GetMemberRole(roomID, updaterID)
	if err != nil {
		return nil, err
	}
	if role != domain.RoleOwner && role != domain.RoleAdmin {
		return nil, errors.New("not authorized to update room")
	}
	room.Name = newName
	if newUsername != "" {
		room.Username = &newUsername
	}
	room.UpdatedAt = time.Now()
	if err := s.roomRepo.Update(room); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *roomService) DeleteRoom(roomID, requesterID string) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil || room == nil {
		return errors.New("room not found")
	}
	role, err := s.membershipRepo.GetMemberRole(roomID, requesterID)
	if err != nil {
		return err
	}
	if role != domain.RoleOwner && role != domain.RoleAdmin {
		return errors.New("not authorized to delete room")
	}
	return s.roomRepo.Delete(roomID)
}

func (s *roomService) AddMember(roomID, requesterID, userID string) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil || room == nil {
		return errors.New("room not found")
	}
	// Check if user is already a member.
	existingRole, err := s.membershipRepo.GetMemberRole(roomID, userID)
	if err != nil {
		return err
	}
	if existingRole != "" {
		return errors.New("user already a member")
	}
	// In channels, only owner/admin can add members.
	if room.Type == domain.RoomTypeChannel {
		reqRole, err := s.membershipRepo.GetMemberRole(roomID, requesterID)
		if err != nil {
			return err
		}
		if reqRole != domain.RoleOwner && reqRole != domain.RoleAdmin {
			return errors.New("not authorized to add member to channel")
		}
	}
	membership := &domain.RoomMembership{
		RoomID:    roomID,
		UserID:    userID,
		Role:      domain.RoleMember,
		CreatedAt: time.Now(),
	}
	return s.membershipRepo.AddMember(membership)
}

func (s *roomService) RemoveMember(roomID, requesterID, userID string) error {
	// If not self-removal, only owner/admin can remove a member.
	if requesterID != userID {
		reqRole, err := s.membershipRepo.GetMemberRole(roomID, requesterID)
		if err != nil {
			return err
		}
		if reqRole != domain.RoleOwner && reqRole != domain.RoleAdmin {
			return errors.New("not authorized to remove member")
		}
	}
	return s.membershipRepo.RemoveMember(roomID, userID)
}

func (s *roomService) PromoteMember(roomID, requesterID, userID string) error {
	reqRole, err := s.membershipRepo.GetMemberRole(roomID, requesterID)
	if err != nil {
		return err
	}
	if reqRole != domain.RoleOwner {
		return errors.New("only owner can promote member")
	}
	return s.membershipRepo.UpdateMemberRole(roomID, userID, domain.RoleAdmin)
}

func (s *roomService) BanMember(roomID, requesterID, userID string) error {
	reqRole, err := s.membershipRepo.GetMemberRole(roomID, requesterID)
	if err != nil {
		return err
	}
	if reqRole != domain.RoleOwner && reqRole != domain.RoleAdmin {
		return errors.New("not authorized to ban member")
	}
	return s.membershipRepo.UpdateMemberRole(roomID, userID, domain.RoleBanned)
}

func (s *roomService) UnbanMember(roomID, requesterID, userID string) error {
	reqRole, err := s.membershipRepo.GetMemberRole(roomID, requesterID)
	if err != nil {
		return err
	}
	if reqRole != domain.RoleOwner && reqRole != domain.RoleAdmin {
		return errors.New("not authorized to unban member")
	}
	return s.membershipRepo.UpdateMemberRole(roomID, userID, domain.RoleMember)
}

