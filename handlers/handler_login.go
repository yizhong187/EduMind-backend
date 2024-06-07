package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yizhong187/EduMind-backend/contextKeys"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/domain"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

// HandlerLogin handles the request to login to an existing user. A cookie containing the JWT will be returned.
func HandlerLogin(w http.ResponseWriter, r *http.Request) {

	apiCfg := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)

	// Decode the JSON request body into parameters struct
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	// Check for empty name or description
	if params.Username == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Username is required")
		return
	} else if params.Password == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	userType, err := apiCfg.DB.GetUserTypeByUsername(r.Context(), params.Username)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't find user type")
		return
	}

	var passwordHash string
	if userType == "tutor" {
		passwordHash, err = apiCfg.DB.GetTutorHash(r.Context(), params.Username)
	} else if userType == "student" {
		passwordHash, err = apiCfg.DB.GetStudentHash(r.Context(), params.Username)
	} else {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid user type")
	}
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving passwordHash: %v", err))
		return
	}

	// Check if the password matches the hashed password in the database
	if !util.CheckPasswordHash(params.Password, passwordHash) {
		util.RespondWithError(w, http.StatusBadRequest, "Wrong password")
		return
	}

	// Query for user
	user, err := apiCfg.DB.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving user info: %v", err))
		return
	}

	// Define the standard claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "github.com/yizhong187/EduMind-backend",
		Subject:   user.UserID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1 day
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(apiCfg.SecretKey))
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Could not login")
		http.Error(w, "could not login", http.StatusInternalServerError)
		return
	}

	// Set the token in an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,                 // Set to false if testing over HTTP
		SameSite: http.SameSiteNoneMode, // Ensure this is set to None for cross-site requests
		Path:     "/",                   // Make sure the cookie is sent with every request to the server
	})

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseUserToUser(user))

}
