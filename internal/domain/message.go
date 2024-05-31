package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
)

type Message struct {
	MessageID uuid.UUID `json:"message_id"`
	ChatID    int32     `json:"chat_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   bool      `json:"deleted"`
	Content   string    `json:"content"`
}

func DatabaseMessageToMessage(message database.Message) Message {
	return Message{
		MessageID: message.MessageID,
		ChatID:    message.ChatID,
		UserID:    message.UserID,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
		Deleted:   message.Deleted,
		Content:   message.Content,
	}
}
