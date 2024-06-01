package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/handlers"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/middlewares"
)

func StudentChatRouter(apiCfg *config.ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", middlewares.MiddlewareStudentAuth(handlers.HandlerStudentGetAllChats, apiCfg))
	r.Post("/", middlewares.MiddlewareStudentAuth(handlers.HandlerStartNewChat, apiCfg))

	// TODO: Add an additional middleware to check if the the studentID from the auth matches the chat's
	//       student_id
	r.Get("/{chatID}", middlewares.MiddlewareStudentAuth(handlers.HandlerStudentGetAllMessages, apiCfg))
	r.Post("/{chatID}", middlewares.MiddlewareStudentAuth(handlers.HandlerStudentNewMessage, apiCfg))

	// r.Delete("/{chatID}/{threadID}", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerDeleteMessage))
	// r.Put("/{chatID}/{threadID}", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerEditMessage))

	//r.Get("/joinChat/{chatID}", middlewares.MiddlewareChatAuth(apiCfg.WSHandler.JoinRoom, apiCfg))
	return r
}
