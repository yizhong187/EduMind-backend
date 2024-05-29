package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/handlers"
)

func StudentRouter(apiCfg *handlers.ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", handlers.HandlerReadiness)
		r.Get("/err", handlers.HandlerError)

		r.Post("/register", apiCfg.HandlerStudentRegistration)

		r.Get("/profile", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerGetStudentProfile))
		r.Put("/profile", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerUpdateStudentProfile))

	})

	return r
}
