package middlewares

import (
	"context"
	"net/http"

	"github.com/yizhong187/EduMind-backend/internal/config"
)

type contextKey string

const (
	ConfigKey contextKey = "config"
)

func MiddlewareApiConfig(next http.HandlerFunc, apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ConfigKey, apiCfg)
		next(w, r.WithContext(ctx))
	}
}
