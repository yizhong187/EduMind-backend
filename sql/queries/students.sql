-- name: CreateNewStudent :one
INSERT INTO students (student_id, username, email, created_at, name, valid, hashed_password)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetStudentById :one
SELECT * FROM students WHERE student_id = $1;

-- name: GetStudentByUsername :one
SELECT * FROM students WHERE username = $1;

-- name: GetStudentHash :one
SELECT hashed_password FROM students WHERE username = $1;

-- name: UpdateStudentProfile :one
UPDATE students SET username = $1, name = $2, email = $3 WHERE student_id = $4
RETURNING *;

-- name: UpdateStudentPassword :exec
UPDATE students SET hashed_password = $1 WHERE student_id = $2;

-- name: StudentGetAllChats :many
SELECT * FROM chats WHERE student_id = $1
ORDER BY created_at DESC;

-- name: StudentCreateNewChat :one
INSERT INTO chats (student_id, created_at, subject_id, header, photo_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;