// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: utils.sql

package database

import (
	"context"
)

const checkEmailTaken = `-- name: CheckEmailTaken :one
SELECT CASE WHEN EXISTS (SELECT 1 FROM users WHERE email = $1) THEN 1 ELSE 0 END
`

func (q *Queries) CheckEmailTaken(ctx context.Context, email string) (int32, error) {
	row := q.db.QueryRowContext(ctx, checkEmailTaken, email)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}

const checkUsernameTaken = `-- name: CheckUsernameTaken :one
SELECT CASE WHEN EXISTS (SELECT 1 FROM users WHERE username = $1) THEN 1 ELSE 0 END
`

func (q *Queries) CheckUsernameTaken(ctx context.Context, username string) (int32, error) {
	row := q.db.QueryRowContext(ctx, checkUsernameTaken, username)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}

const getAllSubjects = `-- name: GetAllSubjects :many
SELECT subject_id, name FROM subjects
`

func (q *Queries) GetAllSubjects(ctx context.Context) ([]Subject, error) {
	rows, err := q.db.QueryContext(ctx, getAllSubjects)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subject
	for rows.Next() {
		var i Subject
		if err := rows.Scan(&i.SubjectID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllTopics = `-- name: GetAllTopics :many
SELECT subject_id, topic_id, name FROM topics
`

func (q *Queries) GetAllTopics(ctx context.Context) ([]Topic, error) {
	rows, err := q.db.QueryContext(ctx, getAllTopics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Topic
	for rows.Next() {
		var i Topic
		if err := rows.Scan(&i.SubjectID, &i.TopicID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllTopicsBySubject = `-- name: GetAllTopicsBySubject :many
SELECT subject_id, topic_id, name FROM topics WHERE subject_id = $1
`

func (q *Queries) GetAllTopicsBySubject(ctx context.Context, subjectID int32) ([]Topic, error) {
	rows, err := q.db.QueryContext(ctx, getAllTopicsBySubject, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Topic
	for rows.Next() {
		var i Topic
		if err := rows.Scan(&i.SubjectID, &i.TopicID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSubjectIDByName = `-- name: GetSubjectIDByName :one
SELECT subject_id
FROM subjects
WHERE name = $1
`

func (q *Queries) GetSubjectIDByName(ctx context.Context, name string) (int32, error) {
	row := q.db.QueryRowContext(ctx, getSubjectIDByName, name)
	var subject_id int32
	err := row.Scan(&subject_id)
	return subject_id, err
}
