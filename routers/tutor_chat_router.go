package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/internal/config"
)

func TutorChatRouter(apiCfg *config.ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	// r.Get("/", apiCfg.MiddlewareTutorAuth(apiCfg.HandlerGetAllChats))
	// r.Post("/", apiCfg.MiddlewareTutorAuth(apiCfg.HandlerStartNewChat))

	// r.Get("/{chatID}", apiCfg.MiddlewareTutorAuth(apiCfg.HandlerChat))
	// r.Post("/{chatID}", apiCfg.MiddlewareTutorAuth(apiCfg.HandlerSendMessage))

	// r.Delete("/{chatID}/{threadID}", apiCfg.MiddlewareTutorAuth(apiCfg.HandlerDeleteMessage))
	// r.Put("/{chatID}/{threadID}", apiCfg.MiddlewareTutorAuth(apiCfg.HandlerEditMessage))

	return r
}
