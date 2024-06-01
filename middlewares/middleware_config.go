package middlewares

import (
	"context"
	"net/http"

	"github.com/yizhong187/EduMind-backend/contextKeys"
	"github.com/yizhong187/EduMind-backend/internal/config"
)

func MiddlewareApiConfig(next http.Handler, apiCfg *config.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKeys.ConfigKey, apiCfg)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
