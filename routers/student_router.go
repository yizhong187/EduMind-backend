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

		rAuthenticated := chi.NewRouter()
		rAuthenticated.Use(middlewares.MiddlewareStudentAuth)
		rAuthenticated.Get("/profile", handlers.HandlerGetStudentProfile)
		rAuthenticated.Put("/profile", handlers.HandlerUpdateStudentProfile)

		r.Mount("/", rAuthenticated)

	})

	return r
}
