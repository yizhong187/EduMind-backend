package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

type Chat struct {
	ChatID    int32      `json:"chat_id"`
	StudentID uuid.UUID  `json:"student_id"`
	TutorID   *uuid.UUID `json:"tutor_id"`
	CreatedAt time.Time  `json:"created_at"`
	SubjectID int32      `json:"subject_id"`
	Topics    []int32    `json:"topics"`
	Header    string     `json:"header"`
	PhotoURL  *string    `json:"photo_url"`
	Completed bool       `json:"completed"`
	Rating    *int32     `json:"rating"`
}

func DatabaseChatToChat(chat database.Chat, topics []int32) Chat {
	return Chat{
		ChatID:    chat.ChatID,
		StudentID: chat.StudentID,
		TutorID:   util.NullUUIDToUUID(chat.TutorID),
		CreatedAt: chat.CreatedAt,
		SubjectID: chat.SubjectID,
		Topics:    topics,
		Header:    chat.Header,
		PhotoURL:  util.NullStringToString(chat.PhotoUrl),
		Completed: chat.Completed,
		Rating:    util.NullInt32ToInt32(chat.Rating),
	}
}
