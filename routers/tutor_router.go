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

		r.Post("/register", handlers.HandlerTutorRegistration)
		r.Post("/login", handlers.HandlerTutorLogin)

		rAuthenticated := chi.NewRouter()
		rAuthenticated.Use(middlewares.MiddlewareTutorAuth)
		rAuthenticated.Get("/profile", handlers.HandlerGetTutorProfile)
		rAuthenticated.Put("/profile", handlers.HandlerUpdateTutorProfile)
		rAuthenticated.Get("/student-profile", handlers.HandlerTutorGetStudentProfile)
		rAuthenticated.Get("/pending-questions", handlers.HandlerGetAvailableQuestions)
		rAuthenticated.Post("/accept", handlers.HandlerAcceptQuestion)

		r.Mount("/", rAuthenticated)
	})

	return r
}
