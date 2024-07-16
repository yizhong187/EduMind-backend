package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
)

type Chat struct {
	ChatID    int32      `json:"chat_id"`
	StudentID uuid.UUID  `json:"student_id"`
	TutorID   *uuid.UUID `json:"tutor_id"`
	CreatedAt time.Time  `json:"created_at"`
	SubjectID int32      `json:"subject"`
	Topics    []int32    `json:"topic"`
	Header    string     `json:"header"`
	PhotoURL  *string    `json:"photo_url"`
	Completed bool       `json:"completed"`
}

func DatabaseChatToChat(chat database.Chat, topics []int32) Chat {
	return Chat{
		ChatID:    chat.ChatID,
		StudentID: chat.StudentID,
		TutorID:   nullUUIDToUUID(chat.TutorID),
		CreatedAt: chat.CreatedAt,
		SubjectID: chat.SubjectID,
		Topics:    topics,
		Header:    chat.Header,
		PhotoURL:  nullStringToString(chat.PhotoUrl),
		Completed: chat.Completed,
	}
}

func nullUUIDToUUID(nu uuid.NullUUID) *uuid.UUID {
	if nu.Valid {
		return &nu.UUID
	}
	return nil
}

func nullStringToString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
