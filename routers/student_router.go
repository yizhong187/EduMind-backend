package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/handlers"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/middlewares"
)

func StudentRouter(apiCfg *config.ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", handlers.HandlerReadiness)
		r.Get("/err", handlers.HandlerError)

		r.Post("/register", (handlers.HandlerStudentRegistration))

		r.Get("/profile", middlewares.MiddlewareStudentAuth(handlers.HandlerGetStudentProfile, apiCfg))
		r.Put("/profile", middlewares.MiddlewareStudentAuth(handlers.HandlerUpdateStudentProfile, apiCfg))

		r.Mount("/{studentID}", StudentChatRouter(apiCfg))

	})

	return r
}
