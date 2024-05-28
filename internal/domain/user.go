package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Valid     bool      `json:"valid"`
}

func DatabaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		Name:      user.Name,
		Valid:     user.Valid,
	}
}
