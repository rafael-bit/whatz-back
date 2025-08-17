package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Avatar    string    `json:"avatar" db:"avatar"`
	Status    string    `json:"status" db:"status"`
	Role      string    `json:"role" db:"role"` // admin, user
	Tags      string    `json:"tags" db:"tags"` // JSON array de tags
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser(username, email, avatar string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Username:  username,
		Email:     email,
		Avatar:    avatar,
		Status:    "online",
		Role:      "user",
		Tags:      "[]", // Array vazio de tags
		CreatedAt: now,
		UpdatedAt: now,
	}
}
