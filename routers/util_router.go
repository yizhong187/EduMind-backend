package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/handlers"
	"github.com/yizhong187/EduMind-backend/internal/config"
)

func UtilRouter(apiCfg *config.ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {

		r.Get("/topics", handlers.HandlerGetAllTopics)

		rSubjects := chi.NewRouter()
		rSubjects.Get("/", handlers.HandlerGetAllSubjects)
		rSubjects.Get("/{subjectID}/topics", handlers.HandlerGetSubjectTopics)

		r.Mount("/subjects", rSubjects)
	})

	return r
}
