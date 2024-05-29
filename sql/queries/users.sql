-- name: InsertNewUser :exec
INSERT INTO users (user_id, username, user_type) VALUES ($1, $2, $3);

-- name: GetUserTypeById :one
SELECT user_type FROM users WHERE user_id = $1;

-- name: GetUserTypeByUsername :one
SELECT user_type FROM users WHERE username = $1;

-- name: CheckUsernameTaken :one
SELECT CASE WHEN EXISTS (SELECT 1 FROM users WHERE username = $1) THEN 1 ELSE 0 END;

-- name: UpdateUsername :exec
UPDATE users SET username = $1 WHERE user_id = $2;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;