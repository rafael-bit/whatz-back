package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"`
	AccessTags  string    `json:"access_tags" db:"access_tags"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewRoom(name, description, roomType, createdBy string) *Room {
	now := time.Now()
	return &Room{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Type:        roomType,
		AccessTags:  "[]",
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
