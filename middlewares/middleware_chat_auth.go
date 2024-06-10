package middlewares

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/contextKeys"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/domain"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

func MiddlewareChatAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
		if !ok || apiCfg == nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
			return
		}

		user, ok := r.Context().Value(contextKeys.UserKey).(domain.User)
		if !ok {
			util.RespondWithError(w, http.StatusInternalServerError, "User not found")
			return
		}

		chatID := chi.URLParam(r, "chatID")
		parsedChatID, err := strconv.ParseInt(chatID, 10, 32)
		fmt.Println("chatID: ", parsedChatID)
		if err != nil {
			fmt.Println(err)
			util.RespondWithError(w, http.StatusBadRequest, "Invalid chat ID")
			return
		}

		chat, err := apiCfg.DB.GetChatById(r.Context(), int32(parsedChatID))
		if err != nil {
			fmt.Println(err)
			util.RespondWithError(w, http.StatusInternalServerError, "Could not get chat details")
			return
		}
		if chat.StudentID != user.ID && (!chat.TutorID.Valid || chat.TutorID.UUID != user.ID) {
			util.RespondWithError(w, http.StatusUnauthorized, "Unauthorized to view chat")
			return
		}

		next.ServeHTTP(w, r)
	})
}
