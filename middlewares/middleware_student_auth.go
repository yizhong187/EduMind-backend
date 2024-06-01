package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

type authedStudentHandler func(http.ResponseWriter, *http.Request, database.Student)

// TODO: REMOVE APICONFIG
func MiddlewareStudentAuth(handler authedStudentHandler, apiCfg *config.ApiConfig) http.HandlerFunc {
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

		userType, err := apiCfg.DB.GetUserTypeById(r.Context(), parsedUUID)
		if userType != "student" {
			util.RespondWithError(w, http.StatusUnauthorized, "Invalid user type")
			return
		}
		if err != nil {
			fmt.Println(err)
			util.RespondWithError(w, http.StatusInternalServerError, "Could not get user type")
			return
		}

		student, err := apiCfg.DB.GetStudentById(r.Context(), parsedUUID)
		if err != nil {
			fmt.Println(err)
			util.RespondWithError(w, http.StatusInternalServerError, "Could not get student details")
			return
		}

		handler(w, r.WithContext(ctx), student)
	}
}
