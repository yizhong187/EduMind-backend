package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/handlers"
)

func StudentChatRouter(apiCfg *handlers.ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerStudentGetAllChats))
	r.Post("/", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerStartNewChat))

	// TODO: Add an additional middleware to check if the the studentID from the auth matches the chat's
	//       student_id
	r.Get("/{chatID}", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerStudentGetAllMessages))
	r.Post("/{chatID}", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerStudentNewMessage))

	// r.Delete("/{chatID}/{threadID}", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerDeleteMessage))
	// r.Put("/{chatID}/{threadID}", apiCfg.MiddlewareStudentAuth(apiCfg.HandlerEditMessage))

	return r
}
