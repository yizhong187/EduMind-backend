-- name: CreateUser :one
INSERT INTO users (id, created_at, name, valid, hashed_password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;