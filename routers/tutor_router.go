package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/handlers"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/middlewares"
)

func TutorRouter(apiCfg *config.ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", handlers.HandlerReadiness)
		r.Get("/err", handlers.HandlerError)

		r.Get("/profile", middlewares.MiddlewareTutorAuth(handlers.HandlerGetTutorProfile, apiCfg))
		r.Put("/profile", middlewares.MiddlewareTutorAuth(handlers.HandlerUpdateTutorProfile, apiCfg))

	})

	return r
}
