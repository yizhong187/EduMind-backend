-- name: CreateNewChat :one
INSERT INTO chats (student_id, created_at, subject_id, header, photo_url)
VALUES ($1, $2, $3, $4, $5)
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

-- name: TutorGetNewChat :one
SELECT * FROM chats WHERE topic IS NULL
ORDER BY created_at ASC
LIMIT 1;

-- name: TutorGetAvailableQuestions :many
SELECT chat_id, student_id, tutor_id, created_at, subject_id, topic, header, photo_url, completed
FROM chats
WHERE topic IS NULL AND subject_id = ANY(
    SELECT ts.subject_id
    FROM tutor_subjects ts
    WHERE ts.tutor_id = $1
);