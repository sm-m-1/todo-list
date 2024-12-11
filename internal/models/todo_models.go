package models

import (
	"time"
)

// Todo struct for the To-Do item, using GORM's model struct
// gorm.Model definition
type Todo struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	UserID      uint   `json:"user_id" gorm:"not null"` // Foreign key to associate with User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User represents a user in the system
// gorm.Model definition
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"` // Hashed password
	CreatedAt time.Time
	UpdatedAt time.Time
	Todos     []Todo `json:"todos" gorm:"foreignKey:UserID"` // One-to-many relationship
}

// Session represents a session in the database for session storage
// gorm.Model definition
type Session struct {
	Token   string    `gorm:"primaryKey;size:255"` // Session token
	Data    []byte    `gorm:"not null"`            // Encoded session data
	Expiry  time.Time `gorm:"not null"`            // Expiry timestamp
	Created time.Time
	Updated time.Time
}
