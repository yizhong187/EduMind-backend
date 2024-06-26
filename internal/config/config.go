package config

import (
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/ws"

	_ "github.com/lib/pq"
)

type ApiConfig struct {
	DB        *database.Queries
	SecretKey string
	WSHandler *ws.Handler
}
