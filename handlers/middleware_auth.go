package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg ApiConfig) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		user, err := apiCfg.DB.GetUserById(r.Context(), parsedUUID)
		if err != nil {
			fmt.Println(err)
			util.RespondWithError(w, http.StatusInternalServerError, "Could not get user info")
			return
		}

		handler(w, r.WithContext(ctx), user)
	}
}
