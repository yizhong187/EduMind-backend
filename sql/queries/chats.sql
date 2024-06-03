-- name: CreateNewChat :one
INSERT INTO chats (student_id, created_at, subject, header)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAllChats :many
SELECT * FROM chats WHERE student_id = $1 OR tutor_id = $1
ORDER BY created_at DESC;

-- name: StudentGetAllChats :many
SELECT * FROM chats WHERE student_id = $1
ORDER BY created_at DESC;

-- name: TutorGetAllChats :many
SELECT * FROM chats WHERE tutor_id = $1
ORDER BY created_at DESC;

-- name: TutorUpdateChat :one
UPDATE chats SET tutor_id = $1, topic = $2 WHERE chat_id = $3
RETURNING *;

-- name: CompleteChat :exec
UPDATE chats SET completed = TRUE WHERE chat_id = $1;

-- name: UpdateChatHeader :one
UPDATE chats SET header = $1 WHERE chat_id = $2
RETURNING *;

-- name: GetChatById :one
SELECT * FROM chats WHERE chat_id = $1;