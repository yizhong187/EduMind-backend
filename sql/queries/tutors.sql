-- name: CreateNewTutor :one
INSERT INTO tutors (tutor_id, username, created_at, name, valid, hashed_password, yoe, subject, verified, rating, rating_count)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetTutorById :one
SELECT * FROM tutors WHERE tutor_id = $1;

-- name: GetTutorByUsername :one
SELECT * FROM tutors WHERE username = $1;

-- name: GetTutorHash :one
SELECT hashed_password FROM tutors WHERE username = $1;