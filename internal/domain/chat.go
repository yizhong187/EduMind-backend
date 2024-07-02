package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
)

type Chat struct {
	ChatID    int32          `json:"chat_id"`
	StudentID uuid.UUID      `json:"student_id"`
	TutorID   uuid.NullUUID  `json:"tutor_id"`
	CreatedAt time.Time      `json:"created_at"`
	Subject   string         `json:"subject"`
	Topic     sql.NullString `json:"topic"`
	Header    string         `json:"header"`
	PhotoURL  sql.NullString `json:"photo_url"`
	Completed bool           `json:"completed"`
}

func DatabaseChatToChat(chat database.Chat) Chat {
	return Chat{
		ChatID:    chat.ChatID,
		StudentID: chat.StudentID,
		TutorID:   chat.TutorID,
		CreatedAt: chat.CreatedAt,
		Subject:   chat.Subject,
		Topic:     chat.Topic,
		Header:    chat.Header,
		PhotoURL:  chat.PhotoUrl,
		Completed: chat.Completed,
	}
}
