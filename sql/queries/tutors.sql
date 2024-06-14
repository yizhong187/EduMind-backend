-- name: CreateNewTutor :one
INSERT INTO tutors (tutor_id, username, email, created_at, name, valid, hashed_password, yoe, subject, verified, rating, rating_count)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetTutorById :one
SELECT * FROM tutors WHERE tutor_id = $1;

-- name: GetTutorByUsername :one
SELECT * FROM tutors WHERE username = $1;

-- name: GetTutorHash :one
SELECT hashed_password FROM tutors WHERE username = $1;

-- name: UpdateTutorProfile :one
UPDATE tutors SET username = $1, name = $2 WHERE tutor_id = $3
RETURNING *;

-- name: UpdateTutorPassword :exec
UPDATE tutors SET hashed_password = $1 WHERE tutor_id = $2;