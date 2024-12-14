package domain

import "time"

// User represents the user entity. We use UUID as primary key.
type User struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"` // UUID string as primary key
	Name      string    `gorm:"not null" json:"name"`
	Phone     string    `gorm:"unique;not null" json:"phone"`
	Username  *string   `gorm:"unique" json:"username,omitempty"` // optional at registration; unique when set
	Password  string    `gorm:"not null" json:"-"`                // hashed password (do not return)
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository defines methods for user persistence.
type UserRepository interface {
	Create(user *User) error
	FindByPhone(phone string) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByID(id string) (*User, error)
	Update(user *User) error
	Delete(user *User) error
}
