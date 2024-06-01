package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

type authedUserHandler func(http.ResponseWriter, *http.Request, database.User)

func MiddlewareUserAuth(handler authedUserHandler, apiCfg *config.ApiConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		chatID := chi.URLParam(r, "chatID")

		cookie, err := r.Cookie("jwt")
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "User unauthenticated")
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(apiCfg.SecretKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				util.RespondWithError(w, http.StatusUnauthorized, "User unauthenticated")
				return
			}
			util.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Bad request: \n%v", err))
			return
		}

		if !token.Valid {
			util.RespondWithError(w, http.StatusUnauthorized, "User unauthenticated")
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			util.RespondWithError(w, http.StatusUnauthorized, "User unauthenticated")
			return
		}

		ctx := context.WithValue(r.Context(), "userClaims", claims)

		parsedUUID, err := uuid.Parse(claims.Subject)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Invalid UUID")
			return
		}

		parsedChatID, err := strconv.ParseInt(chatID, 10, 32)
		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "Invalid chat ID")
			return
		}

		chat, err := apiCfg.DB.GetChatById(r.Context(), int32(parsedChatID))
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Could not get chat details")
			return
		}
		if chat.StudentID != parsedUUID && (!chat.TutorID.Valid || chat.TutorID.UUID != parsedUUID) {
			util.RespondWithError(w, http.StatusUnauthorized, "Unauthorized to view chat")
			return
		}

		user, err := apiCfg.DB.GetUserById(r.Context(), parsedUUID)
		if err != nil {
			fmt.Println(err)
			util.RespondWithError(w, http.StatusInternalServerError, "Could not get user details")
			return
		}

		handler(w, r.WithContext(ctx), user)
	}
}
