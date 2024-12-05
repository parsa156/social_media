package domain

import "time"

// User represents the user entity.
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"uniqueIndex" json:"uuid"`       // Unique code generated for the user.
	Name      string    `json:"name"`                          // Not unique
	Phone     string    `gorm:"uniqueIndex" json:"phone"`        // Unique phone number (used for login)
	Username  string    `gorm:"uniqueIndex" json:"username"`     // Unique username (always starts with '@')
	Password  string    `json:"-"`                             // Hashed password
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	Create(user *User) error
	FindByPhone(phone string) (*User, error)
}
