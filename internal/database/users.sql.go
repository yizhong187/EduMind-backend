// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getUserTypeById = `-- name: GetUserTypeById :one
SELECT user_type FROM users WHERE user_id = $1
`

func (q *Queries) GetUserTypeById(ctx context.Context, userID uuid.UUID) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserTypeById, userID)
	var user_type string
	err := row.Scan(&user_type)
	return user_type, err
}

const getUserTypeByUsername = `-- name: GetUserTypeByUsername :one
SELECT user_type FROM users WHERE username = $1
`

func (q *Queries) GetUserTypeByUsername(ctx context.Context, username string) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserTypeByUsername, username)
	var user_type string
	err := row.Scan(&user_type)
	return user_type, err
}

const insertNewUser = `-- name: InsertNewUser :exec
INSERT INTO users (user_id, username, user_type) VALUES ($1, $2, $3)
`

type InsertNewUserParams struct {
	UserID   uuid.UUID
	Username string
	UserType string
}

func (q *Queries) InsertNewUser(ctx context.Context, arg InsertNewUserParams) error {
	_, err := q.db.ExecContext(ctx, insertNewUser, arg.UserID, arg.Username, arg.UserType)
	return err
}
