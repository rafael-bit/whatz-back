package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string    `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	UserID    string    `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Avatar    string    `json:"avatar" db:"avatar"`
	Type      string    `json:"type" db:"type"` // text, image, file, system
	RoomID    string    `json:"room_id" db:"room_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewMessage(content, userID, username, avatar, messageType, roomID string) *Message {
	now := time.Now().UTC()
	return &Message{
		ID:        uuid.New().String(),
		Content:   content,
		UserID:    userID,
		Username:  username,
		Avatar:    avatar,
		Type:      messageType,
		RoomID:    roomID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
