package config

import (
	"database/sql"

	"github.com/yizhong187/EduMind-backend/internal/database"

	_ "github.com/lib/pq"
)

type ApiConfig struct {
	DB        *database.Queries
	DBConn    *sql.DB
	SecretKey string
}
