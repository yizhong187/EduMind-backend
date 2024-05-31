package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/handlers"
)

func StudentChatRouter(apiCfg *handlers.ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerStudentGetAllChats))
	r.Post("/", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerStartNewChat))

	r.Get("/{chatID}", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerChat))
	r.Post("/{chatID}", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerSendMessage))

	// r.Delete("/{chatID}/{threadID}", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerDeleteMessage))
	// r.Put("/{chatID}/{threadID}", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerEditMessage))

	return r
}
