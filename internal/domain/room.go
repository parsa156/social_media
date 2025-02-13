package domain

import "time"

// RoomType defines the type of room.
type RoomType string

const (
	RoomTypeGroup   RoomType = "group"
	RoomTypeChannel RoomType = "channel"
)

// Room represents a group or channel.
type Room struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Username  *string   `json:"username,omitempty"`
	Type      RoomType  `json:"type"` // "group" or "channel"
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoomMembershipRole defines roles for room membership.
type RoomMembershipRole string

const (
	RoleOwner  RoomMembershipRole = "owner"
	RoleAdmin  RoomMembershipRole = "admin"
	RoleMember RoomMembershipRole = "member"
	RoleBanned RoomMembershipRole = "banned"
)

// RoomMembership represents a user’s membership in a room.
type RoomMembership struct {
	RoomID    string              `json:"room_id"`
	UserID    string              `json:"user_id"`
	Role      RoomMembershipRole  `json:"role"`
	CreatedAt time.Time           `json:"created_at"`
}

// RoomMessage represents a message sent in a room.
type RoomMessage struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository interfaces for room functionality.
type RoomRepository interface {
	Create(room *Room) error
	Update(room *Room) error
	Delete(roomID string) error
	FindByID(roomID string) (*Room, error)
	FindByUsername(username string) (*Room, error)
}

type RoomMembershipRepository interface {
	AddMember(membership *RoomMembership) error
	UpdateMemberRole(roomID, userID string, role RoomMembershipRole) error
	RemoveMember(roomID, userID string) error
	GetMembers(roomID string) ([]*RoomMembership, error)
	IsUserBanned(roomID, userID string) (bool, error)
	GetMemberRole(roomID, userID string) (RoomMembershipRole, error)
}

type RoomMessageRepository interface {
	Create(message *RoomMessage) error
	Update(message *RoomMessage) error
	Delete(messageID string) error
	FindByRoom(roomID string) ([]*RoomMessage, error)
	FindByID(messageID string) (*RoomMessage, error)
}
