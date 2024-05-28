package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/handlers"
)

func ChatRouter(apiCfg *handlers.ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	// Thread-related endpoints
	r.Get("/", handlers.HandlerAllThreads)
	r.Get("/{threadID}", handlers.HandlerThread)

	r.With(util.AuthenticateUserMiddleware).Post("/", handlers.HandlerCreateThread)
	r.With(util.AuthenticateUserMiddleware).Put("/{threadID}", handlers.HandlerUpdateThread)
	r.With(util.AuthenticateUserMiddleware).Delete("/{threadID}", handlers.HandlerDeleteThread)

	// Mount ReplyRouter under a specific thread
	r.Mount("/{threadID}/replies", ReplyRouter())
	return r
}
