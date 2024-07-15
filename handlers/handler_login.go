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

// TO BE DELETED! REPLACED BY STUDENT AND TUTOR SPECIFIC LOGINS
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
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30 days
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

	// Create the response object containing the token and user data
	response := struct {
		Token string      `json:"token"`
		User  domain.User `json:"user"`
	}{
		Token: tokenString,
		User:  domain.DatabaseUserToUser(user),
	}

	// Respond with the token and user data
	util.RespondWithJSON(w, http.StatusOK, response)

}

// HandlerLogin handles the request to login to an existing user. A cookie containing the JWT will be returned.
func HandlerStudentLogin(w http.ResponseWriter, r *http.Request) {

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

	if params.Username == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Username is required")
		return
	} else if params.Password == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	passwordHash, err := apiCfg.DB.GetStudentHash(r.Context(), params.Username)

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

	// Query for student
	student, err := apiCfg.DB.GetStudentByUsername(r.Context(), params.Username)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving student info: %v", err))
		return
	}

	// Define the standard claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "github.com/yizhong187/EduMind-backend",
		Subject:   student.StudentID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30 days
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

	// Create the response object containing the token and user data
	response := struct {
		Token   string         `json:"token"`
		Student domain.Student `json:"student"`
	}{
		Token:   tokenString,
		Student: domain.DatabaseStudentToStudent(student),
	}

	// Respond with the token and user data
	util.RespondWithJSON(w, http.StatusOK, response)

}

// HandlerLogin handles the request to login to an existing user. A cookie containing the JWT will be returned.
func HandlerTutorLogin(w http.ResponseWriter, r *http.Request) {

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

	if params.Username == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Username is required")
		return
	} else if params.Password == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	passwordHash, err := apiCfg.DB.GetTutorHash(r.Context(), params.Username)

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

	// Query for tutor
	tutor, err := apiCfg.DB.GetTutorByUsername(r.Context(), params.Username)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving tutor info: %v", err))
		return
	}

	// Define the standard claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "github.com/yizhong187/EduMind-backend",
		Subject:   tutor.TutorID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30 days
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

	tutorSubjects, err := apiCfg.DB.GetTutorSubjects(r.Context(), tutor.TutorID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get tutor subjects")
		return
	}

	// Create the response object containing the token and user data
	response := struct {
		Token string       `json:"token"`
		Tutor domain.Tutor `json:"tutor"`
	}{
		Token: tokenString,
		Tutor: domain.DatabaseTutorToTutor(tutor, tutorSubjects),
	}

	// Respond with the token and user data
	util.RespondWithJSON(w, http.StatusOK, response)

}
