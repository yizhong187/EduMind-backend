package handlers

import (
	"github.com/yizhong187/EduMind-backend/internal/database"

	_ "github.com/lib/pq"
)

type ApiConfig struct {
	DB        *database.Queries
	SecretKey string
}
