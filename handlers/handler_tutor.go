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

	type parameters struct {
		Username string         `json:"username"`
		Name     string         `json:"name"`
		Email    string         `json:"email"`
		Subjects map[string]int `json:"subjects"`
		Password string         `json:"password"`
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
	} else if params.Name == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Name is required")
		return
	} else if params.Email == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Email is required")
		return
	} else if len(params.Subjects) == 0 {
		util.RespondWithError(w, http.StatusBadRequest, "At least one subject is required")
		return
	} else if params.Password == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	usernameTaken, err := apiCfg.DB.CheckUsernameTaken(r.Context(), params.Username)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't check if username taken")
		return
	} else if usernameTaken == 1 {
		util.RespondWithError(w, http.StatusConflict, "Username taken")
		return
	}

	emailTaken, err := apiCfg.DB.CheckEmailTaken(r.Context(), params.Email)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't check if email taken")
		return
	} else if emailTaken == 1 {
		util.RespondWithError(w, http.StatusConflict, "Email taken")
		return
	}

	hashedPassword, err := util.HashPassword(params.Password)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	tx, err := apiCfg.DBConn.BeginTx(r.Context(), nil)

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't start transaction")
		return
	}
	defer tx.Rollback()

	queries := apiCfg.DB.WithTx(tx)

	tutorID := uuid.New()

	err = queries.InsertNewUser(r.Context(), database.InsertNewUserParams{
		UserID:   tutorID,
		Username: params.Username,
		Email:    params.Email,
		UserType: "tutor",
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't insert into user database")
		return
	}

	tutor, err := queries.CreateNewTutor(r.Context(), database.CreateNewTutorParams{
		TutorID:        tutorID,
		Username:       params.Username,
		Email:          params.Email,
		CreatedAt:      time.Now().UTC(),
		Name:           params.Name,
		Valid:          true,
		HashedPassword: hashedPassword,
		Verified:       false,
		Rating:         sql.NullFloat64{Float64: 0.0, Valid: false},
		RatingCount:    0,
	})
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create tutor")
		return
	}

	for subjectName, yoe := range params.Subjects {
		subjectID, err := queries.GetSubjectIDByName(r.Context(), subjectName)
		if err != nil {
			fmt.Println(err)
			util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get subject: %s", subjectName))
			return
		}

		_, err = queries.AddTutorSubject(r.Context(), database.AddTutorSubjectParams{
			TutorID:   tutorID,
			SubjectID: subjectID,
			Yoe:       int32(yoe),
		})

		if err != nil {
			fmt.Println(err)
			util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create tutor-subject relationship")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't commit transaction")
		return
	}

	tutorSubjects, err := apiCfg.DB.GetTutorSubjects(r.Context(), tutorID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get tutor subjects")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, domain.DatabaseTutorToTutor(tutor, tutorSubjects))
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
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't decode parameters: \n%v", err))
		return
	}
	defer r.Body.Close()

	if params.StudentID == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Student ID is required")
		return
	}

	parsedUUID, err := uuid.Parse(params.StudentID)
	if err != nil {
		fmt.Println(err)
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
		Email    string `json:"email"`
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
	} else if params.Name == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Name is required")
		return
	} else if params.Email == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	usernameTaken, err := apiCfg.DB.CheckUsernameTaken(r.Context(), params.Username)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't check if username taken")
		return
	} else if tutor.Username != params.Username && usernameTaken == 1 {
		util.RespondWithError(w, http.StatusConflict, "Username taken")
		return
	}

	emailTaken, err := apiCfg.DB.CheckEmailTaken(r.Context(), params.Email)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't check if email taken")
		return
	} else if tutor.Email != params.Email && emailTaken == 1 {
		util.RespondWithError(w, http.StatusConflict, "Email taken")
		return
	}

	tx, err := apiCfg.DBConn.BeginTx(r.Context(), nil)

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't start transaction")
		return
	}
	defer tx.Rollback()

	queries := apiCfg.DB.WithTx(tx)

	updatedTutor, err := queries.UpdateTutorProfile(r.Context(), database.UpdateTutorProfileParams{
		TutorID:  tutor.TutorID,
		Username: params.Username,
		Name:     params.Name,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update tutor profile")
		return
	}

	err = queries.UpdateUserProfile(r.Context(), database.UpdateUserProfileParams{
		Username: params.Username,
		UserID:   tutor.TutorID,
		Email:    tutor.Email,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update user profile")
		return
	}

	if err := tx.Commit(); err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't commit transaction")
		return
	}

	tutorSubjects, err := apiCfg.DB.GetTutorSubjects(r.Context(), tutor.TutorID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get tutor subjects")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseTutorToTutor(updatedTutor, tutorSubjects))
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
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	if params.OldPassword == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Old password is required")
		return
	} else if params.NewPassword == "" {
		util.RespondWithError(w, http.StatusBadRequest, "New password is required")
		return
	}

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

// TO BE DELETED
func HandlerGetNewChat(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	chat, err := apiCfg.DB.TutorGetNewChat(r.Context())

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Could not get a new chat")
		return
	}

	if chat.ChatID == 0 {
		util.RespondWithJSON(w, http.StatusOK, struct{}{})
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseChatToChat(chat))
}

func HandlerGetAvailableQuestions(w http.ResponseWriter, r *http.Request) {
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

	databaseChats, err := apiCfg.DB.TutorGetAvailableQuestions(r.Context(), tutor.TutorID)

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Could not get a new chat")
		return
	}

	if len(databaseChats) == 0 {
		util.RespondWithJSON(w, http.StatusOK, struct{}{})
	}

	var chats []domain.Chat
	for _, chat := range databaseChats {
		chats = append(chats, domain.DatabaseChatToChat(chat))
	}

	util.RespondWithJSON(w, http.StatusOK, chats)
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
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Invalid chat ID")
		return
	}

	type parameters struct {
		Topic string `json:"topic"`
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	if params.Topic == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Topic is required")
		return
	}

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
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update chat topic")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, chat)
}
