-- name: CreateNewTutor :one
INSERT INTO tutors (tutor_id, username, email, created_at, name, valid, hashed_password, verified, rating, rating_count)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: AddTutorSubject :one
INSERT INTO tutor_subjects (tutor_id, subject_id, yoe)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTutorSubjects :many
SELECT 
    s.name AS subject,
    ts.yoe
FROM 
    tutor_subjects ts
JOIN 
    subjects s ON ts.subject_id = s.subject_id
WHERE 
    ts.tutor_id = $1;

-- name: GetTutorSubjectIDs :many
SELECT subject_id
FROM tutor_subjects
WHERE tutor_id = $1;

-- name: GetSubjectIDByName :one
SELECT subject_id
FROM subjects
WHERE name = $1;

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