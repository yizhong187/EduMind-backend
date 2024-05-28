-- name: CreateNewStudent :one
INSERT INTO students (student_id, username, created_at, name, valid, hashed_password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetStudentById :one
SELECT * FROM students WHERE student_id = $1;

-- name: GetStudentByUsername :one
SELECT * FROM students WHERE username = $1;

-- name: GetStudentHash :one
SELECT hashed_password FROM students WHERE username = $1;