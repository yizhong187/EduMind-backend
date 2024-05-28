package domain

import (
	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
)

type User struct {
	ID       uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Type     string    `json:"user_type"`
}

func DatabaseUserToUser(user database.User) User {
	return User{
		ID:       user.UserID,
		Username: user.Username,
		Type:     user.UserType,
	}
}
