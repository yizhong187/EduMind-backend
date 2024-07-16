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
	SubjectID int32          `json:"subject"`
	Topics    []int32        `json:"topic"`
	Header    string         `json:"header"`
	PhotoURL  sql.NullString `json:"photo_url"`
	Completed bool           `json:"completed"`
}

func DatabaseChatToChat(chat database.Chat, topics []int32) Chat {
	return Chat{
		ChatID:    chat.ChatID,
		StudentID: chat.StudentID,
		TutorID:   chat.TutorID,
		CreatedAt: chat.CreatedAt,
		SubjectID: chat.SubjectID,
		Topics:    topics,
		Header:    chat.Header,
		PhotoURL:  chat.PhotoUrl,
		Completed: chat.Completed,
	}
}
