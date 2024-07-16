-- name: CheckUsernameTaken :one
SELECT CASE WHEN EXISTS (SELECT 1 FROM users WHERE username = $1) THEN 1 ELSE 0 END;

-- name: CheckEmailTaken :one
SELECT CASE WHEN EXISTS (SELECT 1 FROM users WHERE email = $1) THEN 1 ELSE 0 END;

-- name: GetSubjectIDByName :one
SELECT subject_id
FROM subjects
WHERE name = $1;

-- name: GetAllSubjects :many
SELECT * FROM subjects;

-- name: GetAllTopicsBySubject :many
SELECT * FROM topics WHERE subject_id = $1;

-- name: GetAllTopics :many
SELECT * FROM topics;
