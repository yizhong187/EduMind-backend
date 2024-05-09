package handlers

import (
	"github.com/yizhong187/EduMind-backend/internal/database"

	_ "github.com/lib/pq"
)

type ApiConfig struct {
	DB *database.Queries
}

func LoadConfiguration() (*ApiConfig, error) {
	// Load your configuration from a file, environment variables, etc.
	return &ApiConfig{}, nil
}
