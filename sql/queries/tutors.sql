-- name: CreateNewTutor :one
INSERT INTO tutors (tutor_id, username, email, created_at, name, valid, hashed_password, verified, rating, rating_count, photo_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: AddTutorSubject :one
INSERT INTO tutor_subjects (tutor_id, subject_id, yoe)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTutorSubjects :many
SELECT subject_id, yoe FROM tutor_subjects WHERE tutor_id = $1;

-- name: GetTutorSubjectIDs :many
SELECT subject_id FROM tutor_subjects WHERE tutor_id = $1;

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

-- name: TutorGetAvailableQuestions :many
SELECT chat_id, student_id, tutor_id, created_at, subject_id, topic, header, photo_url, completed
FROM chats
WHERE tutor_id IS NULL AND subject_id = ANY(
    SELECT ts.subject_id
    FROM tutor_subjects ts
    WHERE ts.tutor_id = $1
);

-- name: TutorGetAllChats :many
SELECT * FROM chats WHERE tutor_id = $1
ORDER BY created_at DESC;

-- name: TutorUpdateChat :one
UPDATE chats SET tutor_id = $1, topic = $2 WHERE chat_id = $3
RETURNING *;

-- name: TutorAcceptQuestion :exec
UPDATE chats SET tutor_id = $1 WHERE chat_id = $2;