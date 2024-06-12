package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/contextKeys"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/domain"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

func MiddlewareStudentAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
		if !ok || apiCfg == nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
			return
		}

		// Extract the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			util.RespondWithError(w, http.StatusUnauthorized, "Authorization header missing")
			return
		}

		// Check if the header starts with "Bearer "
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			util.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		// Extract the token from the header
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(apiCfg.SecretKey), nil
		})

		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
			util.RespondWithError(w, http.StatusInternalServerError, "Invalid UUID")
			return
		}

		userType, err := apiCfg.DB.GetUserTypeById(ctx, parsedUUID)
		if userType != "student" {
			util.RespondWithError(w, http.StatusUnauthorized, "Invalid user type")
			return
		}
		if err != nil {
			fmt.Println(err)
			util.RespondWithError(w, http.StatusInternalServerError, "Could not get user type")
			return
		}

		student, err := apiCfg.DB.GetStudentById(ctx, parsedUUID)
		if err != nil {
			fmt.Println(err)
			util.RespondWithError(w, http.StatusInternalServerError, "Could not get student details")
			return
		}

		ctx = context.WithValue(ctx, contextKeys.StudentKey, domain.DatabaseStudentToStudent(student))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
