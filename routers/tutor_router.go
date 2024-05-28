package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/handlers"
)

func TutorRouter(apiCfg *handlers.ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", handlers.HandlerReadiness)
		r.Get("/err", handlers.HandlerError)

		r.Get("/profile", apiCfg.MiddlewareTutorAuth(apiCfg.HandlerGetTutorProfile))
		// r.Put("/profile", apiCfg.MiddlewareTutorAuth(apiCfg.HandlerUpdateTutorProfile))

	})

	return r
}
