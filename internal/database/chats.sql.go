// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: chats.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const completeChat = `-- name: CompleteChat :exec
UPDATE chats SET completed = TRUE WHERE chat_id = $1
`

func (q *Queries) CompleteChat(ctx context.Context, chatID int32) error {
	_, err := q.db.ExecContext(ctx, completeChat, chatID)
	return err
}

const createNewChat = `-- name: CreateNewChat :one
INSERT INTO chats (student_id, created_at, subject, header)
VALUES ($1, $2, $3, $4)
RETURNING chat_id, student_id, tutor_id, created_at, subject, topic, header, completed
`

type CreateNewChatParams struct {
	StudentID uuid.UUID
	CreatedAt time.Time
	Subject   string
	Header    string
}

func (q *Queries) CreateNewChat(ctx context.Context, arg CreateNewChatParams) (Chat, error) {
	row := q.db.QueryRowContext(ctx, createNewChat,
		arg.StudentID,
		arg.CreatedAt,
		arg.Subject,
		arg.Header,
	)
	var i Chat
	err := row.Scan(
		&i.ChatID,
		&i.StudentID,
		&i.TutorID,
		&i.CreatedAt,
		&i.Subject,
		&i.Topic,
		&i.Header,
		&i.Completed,
	)
	return i, err
}

const studentGetAllChats = `-- name: StudentGetAllChats :many
SELECT chat_id, student_id, tutor_id, created_at, subject, topic, header, completed FROM chats WHERE student_id = $1
ORDER BY created_at DESC
`

func (q *Queries) StudentGetAllChats(ctx context.Context, studentID uuid.UUID) ([]Chat, error) {
	rows, err := q.db.QueryContext(ctx, studentGetAllChats, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chat
	for rows.Next() {
		var i Chat
		if err := rows.Scan(
			&i.ChatID,
			&i.StudentID,
			&i.TutorID,
			&i.CreatedAt,
			&i.Subject,
			&i.Topic,
			&i.Header,
			&i.Completed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const tutorGetAllChats = `-- name: TutorGetAllChats :many
SELECT chat_id, student_id, tutor_id, created_at, subject, topic, header, completed FROM chats WHERE tutor_id = $1
ORDER BY created_at DESC
`

func (q *Queries) TutorGetAllChats(ctx context.Context, tutorID uuid.NullUUID) ([]Chat, error) {
	rows, err := q.db.QueryContext(ctx, tutorGetAllChats, tutorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chat
	for rows.Next() {
		var i Chat
		if err := rows.Scan(
			&i.ChatID,
			&i.StudentID,
			&i.TutorID,
			&i.CreatedAt,
			&i.Subject,
			&i.Topic,
			&i.Header,
			&i.Completed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const tutorUpdateChat = `-- name: TutorUpdateChat :one
UPDATE chats SET tutor_id = $1, topic = $2 WHERE chat_id = $3
RETURNING chat_id, student_id, tutor_id, created_at, subject, topic, header, completed
`

type TutorUpdateChatParams struct {
	TutorID uuid.NullUUID
	Topic   sql.NullString
	ChatID  int32
}

func (q *Queries) TutorUpdateChat(ctx context.Context, arg TutorUpdateChatParams) (Chat, error) {
	row := q.db.QueryRowContext(ctx, tutorUpdateChat, arg.TutorID, arg.Topic, arg.ChatID)
	var i Chat
	err := row.Scan(
		&i.ChatID,
		&i.StudentID,
		&i.TutorID,
		&i.CreatedAt,
		&i.Subject,
		&i.Topic,
		&i.Header,
		&i.Completed,
	)
	return i, err
}

const updateChatHeader = `-- name: UpdateChatHeader :one
UPDATE chats SET header = $1 WHERE chat_id = $2
RETURNING chat_id, student_id, tutor_id, created_at, subject, topic, header, completed
`

type UpdateChatHeaderParams struct {
	Header string
	ChatID int32
}

func (q *Queries) UpdateChatHeader(ctx context.Context, arg UpdateChatHeaderParams) (Chat, error) {
	row := q.db.QueryRowContext(ctx, updateChatHeader, arg.Header, arg.ChatID)
	var i Chat
	err := row.Scan(
		&i.ChatID,
		&i.StudentID,
		&i.TutorID,
		&i.CreatedAt,
		&i.Subject,
		&i.Topic,
		&i.Header,
		&i.Completed,
	)
	return i, err
}
