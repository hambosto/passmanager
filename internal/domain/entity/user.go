package entity

import "time"

// User represents a vault user (for future multi-user support)
type User struct {
	ID        string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user
func NewUser(email string) *User {
	now := time.Now()
	return &User{
		ID:        generateID(),
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
