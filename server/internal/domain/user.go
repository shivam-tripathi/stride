package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new User
func NewUser(name, email string) *User {
	now := time.Now()
	return &User{
		ID:        GenerateID(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GenerateID generates a new ID for a user
// In a real application, you might use UUID or another ID generation strategy
func GenerateID() string {
	return time.Now().Format("20060102150405") + "-user"
}
