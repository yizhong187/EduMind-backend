-- name: CreateUser :one
INSERT INTO users (id, username, created_at, name, valid, hashed_password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetHash :one
SELECT hashed_password FROM users WHERE username = $1;
