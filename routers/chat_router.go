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

	r.With(middlewares.MiddlewareTutorAuth).Get("/pending", handlers.HandlerGetAvailableQuestions)
	r.With(middlewares.MiddlewareTutorAuth).Post("/{chatID}/accept", handlers.HandlerAcceptQuestion)
	r.With(middlewares.MiddlewareTutorAuth).With(middlewares.MiddlewareChatAuth).Put("/{chatID}/update-topics", handlers.HandlerUpdateChatTopics)

	r.With(middlewares.MiddlewareChatAuth).Get("/{chatID}", handlers.HandlerGetAllMessages)
	r.With(middlewares.MiddlewareChatAuth).Post("/{chatID}/new", handlers.HandlerNewMessage)

	return r
}
