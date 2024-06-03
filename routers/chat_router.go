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

	r.With(middlewares.MiddlewareStudentAuth).Post("/new", handlers.HandlerStartNewChat)
	r.Get("/", handlers.HandlerGetAllChats)

	rAuthenticated := chi.NewRouter()
	rAuthenticated.Use(middlewares.MiddlewareChatAuth)

	rAuthenticated.Get("/{chatID}", handlers.HandlerGetAllMessages)
	rAuthenticated.Get("/joinChat/{chatID}", apiCfg.WSHandler.JoinRoom)

	r.Mount("/", rAuthenticated)

	return r
}
