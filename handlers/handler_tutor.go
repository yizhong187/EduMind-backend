package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/contextKeys"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/domain"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

func HandlerTutorRegistration(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	// local struct to hold expected data from the request body
	type parameters struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		YOE      int    `json:"yoe"`
		Subject  string `json:"subject"`
		Password string `json:"password"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
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
		log.Fatal(err)
	}

	tutor, err := apiCfg.DB.CreateNewTutor(r.Context(), database.CreateNewTutorParams{
		TutorID:        uuid.New(),
		Username:       params.Username,
		CreatedAt:      time.Now().UTC(),
		Name:           params.Name,
		Valid:          true,
		HashedPassword: hashedPassword,
		Yoe:            int32(params.YOE),
		Subject:        params.Subject,
		Verified:       false,
		RatingCount:    0,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, domain.DatabaseTutorToTutor(tutor))
}

func HandlerTutorGetStudentProfile(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	type parameters struct {
		StudentID string `json:"student_id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't decode parameters: \n%v", err))
		return
	}
	defer r.Body.Close()

	parsedUUID, err := uuid.Parse(params.StudentID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Invalid UUID")
		return
	}

	student, err := apiCfg.DB.GetStudentById(r.Context(), parsedUUID)

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Could not get user info")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseStudentToStudent(student))
}

func HandlerUpdateTutorProfile(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	tutor, ok := r.Context().Value(contextKeys.TutorKey).(domain.Tutor)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Tutor profile not found")
		return
	}

	type parameters struct {
		Username string `json:"username"`
		Name     string `json:"name"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
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

	updateTutor, err := apiCfg.DB.UpdateTutorProfile(r.Context(), database.UpdateTutorProfileParams{
		TutorID:  tutor.TutorID,
		Username: params.Username,
		Name:     params.Name,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update tutor profile")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseTutorToTutor(updateTutor))
}

func HandlerUpdateTutorPassword(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	tutor, ok := r.Context().Value(contextKeys.TutorKey).(domain.Tutor)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Tutor profile not found")
		return
	}

	type parameters struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
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

	databaseHashedPassword, err := apiCfg.DB.GetTutorHash(r.Context(), tutor.Username)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve hash")
		return
	}

	passwordMatched := util.CheckPasswordHash(hashedOldPassword, databaseHashedPassword)
	if passwordMatched == false {
		util.RespondWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	hashedNewPassword, err := util.HashPassword(params.NewPassword)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Hashing password went wrong")
		return
	}

	err = apiCfg.DB.UpdateTutorPassword(r.Context(), database.UpdateTutorPasswordParams{
		TutorID:        tutor.TutorID,
		HashedPassword: hashedNewPassword,
	})
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update password")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, struct{}{})
}

func HandlerGetTutorProfile(w http.ResponseWriter, r *http.Request) {
	tutor, ok := r.Context().Value(contextKeys.TutorKey).(domain.Tutor)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Tutor profile not found")
		return
	}
	util.RespondWithJSON(w, http.StatusOK, tutor)
}

func HandlerConfigNewChat(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	tutor, ok := r.Context().Value(contextKeys.TutorKey).(domain.Tutor)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Tutor profile not found")
		return
	}

	chatIDString := chi.URLParam(r, "chatID")
	chatID, err := strconv.ParseInt(chatIDString, 10, 32)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Invalid chat ID")
		return
	}

	type parameters struct {
		Topic string `json:"topic"`
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	chat, err := apiCfg.DB.TutorUpdateChat(r.Context(), database.TutorUpdateChatParams{
		TutorID: uuid.NullUUID{
			UUID:  tutor.TutorID,
			Valid: tutor.TutorID != uuid.Nil,
		},
		Topic: sql.NullString{
			String: params.Topic,
			Valid:  params.Topic != "",
		},
		ChatID: int32(chatID),
	})

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update chat topic")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, chat)
}
