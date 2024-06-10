package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/contextKeys"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/domain"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

func HandlerStudentRegistration(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	// local struct to hold expected data from the request body
	type parameters struct {
		Username string `json: "username"`
		Password string `json: "password"`
		Name     string `json: "name"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	usernameTaken, err := apiCfg.DB.CheckUsernameTaken(r.Context(), params.Username)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't check if username taken")
		return
	} else if usernameTaken == 1 {
		util.RespondWithError(w, http.StatusConflict, "Username taken")
		return
	}

	hashedPassword, err := util.HashPassword(params.Password)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Hashing password went wrong")
		return
	}

	studentUUID := uuid.New()

	err = apiCfg.DB.InsertNewUser(r.Context(), database.InsertNewUserParams{
		UserID:   studentUUID,
		Username: params.Username,
		UserType: "student",
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't insert into user database")
		return
	}

	student, err := apiCfg.DB.CreateNewStudent(r.Context(), database.CreateNewStudentParams{
		StudentID:      studentUUID,
		CreatedAt:      time.Now().UTC(),
		Username:       params.Username,
		Name:           params.Name,
		Valid:          true,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create new student")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, domain.DatabaseStudentToStudent(student))
}

func HandlerUpdateStudentProfile(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	student, ok := r.Context().Value(contextKeys.StudentKey).(domain.Student)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Student profile not found")
		return
	}

	type parameters struct {
		Username string `json: "username"`
		Name     string `json: "name"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	usernameTaken, err := apiCfg.DB.CheckUsernameTaken(r.Context(), params.Username)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't check if username taken")
		return
	} else if student.Username != params.Username && usernameTaken == 1 {
		util.RespondWithError(w, http.StatusConflict, "Username taken")
		return
	}

	updatedStudent, err := apiCfg.DB.UpdateStudentProfile(r.Context(), database.UpdateStudentProfileParams{
		StudentID: student.StudentID,
		Username:  params.Username,
		Name:      params.Name,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update student profile")
		return
	}

	err = apiCfg.DB.UpdateUsername(r.Context(), database.UpdateUsernameParams{
		Username: params.Username,
		UserID:   student.StudentID,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update user profile")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseStudentToStudent(updatedStudent))
}

func HandlerUpdateStudentPassword(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	student, ok := r.Context().Value(contextKeys.StudentKey).(domain.Student)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Student profile not found")
		return
	}

	type parameters struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	hashedOldPassword, err := util.HashPassword(params.OldPassword)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Hashing password went wrong")
		return
	}

	databaseHashedPassword, err := apiCfg.DB.GetTutorHash(r.Context(), student.Username)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve hash")
		return
	}

	passwordMatched := util.CheckPasswordHash(hashedOldPassword, databaseHashedPassword)
	if !passwordMatched {
		util.RespondWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	hashedNewPassword, err := util.HashPassword(params.NewPassword)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Hashing password went wrong")
		return
	}

	err = apiCfg.DB.UpdateStudentPassword(r.Context(), database.UpdateStudentPasswordParams{
		StudentID:      student.StudentID,
		HashedPassword: hashedNewPassword,
	})
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update password")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, struct{}{})
}

func HandlerGetStudentProfile(w http.ResponseWriter, r *http.Request) {
	student, ok := r.Context().Value(contextKeys.StudentKey).(domain.Student)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}
	util.RespondWithJSON(w, http.StatusOK, student)
}

func HandlerStartNewChat(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	student, ok := r.Context().Value(contextKeys.StudentKey).(domain.Student)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	// local struct to hold expected data from the request body
	type parameters struct {
		Subject string `json:"subject"`
		Header  string `json:"header"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	chat, err := apiCfg.DB.CreateNewChat(r.Context(), database.CreateNewChatParams{
		StudentID: student.StudentID,
		CreatedAt: time.Now().UTC(),
		Subject:   params.Subject,
		Header:    params.Header,
	})
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create new chat")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, domain.DatabaseChatToChat(chat))
}
