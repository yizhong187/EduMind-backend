-- name: InsertNewUser :exec
INSERT INTO users (user_id, username, user_type) VALUES ($1, $2, $3);

-- name: GetUserTypeById :one
SELECT user_type FROM users WHERE user_id = $1;

-- name: GetUserTypeByUsername :one
SELECT user_type FROM users WHERE username = $1;
