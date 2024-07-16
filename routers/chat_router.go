package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/handlers"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/middlewares"
)

func ChatRouter(apiCfg *config.ApiConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.MiddlewareUserAuth)

	r.Get("/", handlers.HandlerGetAllChats)

	rAuthenticated := chi.NewRouter()
	rAuthenticated.Route("/{chatID}", func(r chi.Router) {
		r.Use(middlewares.MiddlewareChatAuth)
		r.Get("/view", handlers.HandlerGetAllMessages)
		r.Post("/new", handlers.HandlerNewMessage)
		r.With(middlewares.MiddlewareTutorAuth).Put("/update-topics", handlers.HandlerUpdateChatTopics)
	})

	r.Mount("/", rAuthenticated)

	return r
}
