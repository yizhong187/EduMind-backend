-- name: GetAllMessages :many
SELECT * FROM messages WHERE chat_id = $1
ORDER BY created_at DESC;

-- name: CreateNewMessage :exec
INSERT INTO messages (message_id, chat_id, user_id, created_at, content)
VALUES ($1, $2, $3, $4, $5);

-- name: EditMessage :exec
UPDATE messages SET content = $1, updated_at = $2 WHERE message_id = $3;

-- name: DeleteMessage :exec
UPDATE messages SET deleted = TRUE WHERE message_id = $1;




