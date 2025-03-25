package message

import (
	"time"
)

// Message represents a message in the system
type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewMessage creates a new message
func NewMessage(userID, content string) *Message {
	now := time.Now()
	return &Message{
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}