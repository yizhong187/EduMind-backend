-- name: CreateUser :one
INSERT INTO users (id, created_at, name, valid)
VALUES ($1, $2, $3, $4)
RETURNING *;