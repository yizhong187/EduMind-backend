-- name: GetAllChats :many
SELECT * FROM chats WHERE student_id = $1 OR tutor_id = $1
ORDER BY created_at DESC;

-- name: CompleteChat :exec
UPDATE chats SET completed = TRUE WHERE chat_id = $1;

-- name: UpdateChatHeader :one
UPDATE chats SET header = $1 WHERE chat_id = $2
RETURNING *;

-- name: GetChatById :one
SELECT * FROM chats WHERE chat_id = $1;

-- name: AddChatTopic :exec
INSERT INTO chat_topics (chat_id, topic_id)
VALUES ($1, $2);

-- name: GetChatTopics :many
SELECT topic_id FROM chat_topics WHERE chat_id = $1;

-- name: CheckChatTakenOrCompleted :one
SELECT 
    CASE 
        WHEN tutor_id IS NOT NULL OR completed = TRUE THEN 1 
        ELSE 0 
    END AS is_chat_taken_or_completed
FROM chats
WHERE chat_id = $1;

-- name: UpdateChatRating :one
UPDATE chats SET rating = $1 WHERE chat_id = $2 RETURNING *;