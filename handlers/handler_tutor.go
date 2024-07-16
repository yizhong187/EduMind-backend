package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	type parameters struct {
		Username string      `json:"username"`
		Name     string      `json:"name"`
		Email    string      `json:"email"`
		Subjects map[int]int `json:"subjects"`
		Password string      `json:"password"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Couldn't decode parameters: ", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.Username == "" || params.Name == "" || params.Email == "" || len(params.Subjects) == 0 || params.Password == "" {
		fmt.Println("Missing one more more required parameters.")
		util.RespondWithMissingParametersError(w)
		return
	}

	usernameTaken, err := apiCfg.DB.CheckUsernameTaken(r.Context(), params.Username)
	if err != nil {
		fmt.Println("Couldn't check if username taken: ", err)
		util.RespondWithInternalServerError(w)
		return
	} else if usernameTaken == 1 {
		fmt.Println("Username submitted clashes with existing username.")
		util.RespondWithError(w, http.StatusConflict, "Username already taken")
		return
	}

	emailTaken, err := apiCfg.DB.CheckEmailTaken(r.Context(), params.Email)
	if err != nil {
		fmt.Println("Couldn't check if email taken: ", err)
		util.RespondWithInternalServerError(w)
		return
	} else if emailTaken == 1 {
		fmt.Println("Email submitted clashes with existing email.")
		util.RespondWithError(w, http.StatusConflict, "Email already taken")
		return
	}

	hashedPassword, err := util.HashPassword(params.Password)
	if err != nil {
		fmt.Println("Hashing password went wrong: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	tx, err := apiCfg.DBConn.BeginTx(r.Context(), nil)
	if err != nil {
		fmt.Println("Couldn't start transaction: ", err)
		util.RespondWithInternalServerError(w)
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
		fmt.Println("Couldn't insert into user database: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	_, err = queries.CreateNewTutor(r.Context(), database.CreateNewTutorParams{
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
		fmt.Println("Couldn't create new tutor: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	for subjectID, yoe := range params.Subjects {

		_, err = queries.AddTutorSubject(r.Context(), database.AddTutorSubjectParams{
			TutorID:   tutorID,
			SubjectID: int32(subjectID),
			Yoe:       int32(yoe),
		})

		if err != nil {
			fmt.Println("Couldn't create tutor-subject relationship", err)
			util.RespondWithInternalServerError(w)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Couldn't commit transaction: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, "Registration successful.")
}

func HandlerTutorGetStudentProfile(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	type parameters struct {
		StudentID string `json:"student_id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Couldn't decode parameters:", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.StudentID == "" {
		fmt.Println("Missing student ID.")
		util.RespondWithMissingParametersError(w)
		return
	}

	parsedUUID, err := uuid.Parse(params.StudentID)
	if err != nil {
		fmt.Println("Invalid UUID", err)
		util.RespondWithInternalServerError(w)
		return
	}

	student, err := apiCfg.DB.GetStudentById(r.Context(), parsedUUID)

	if err != nil {
		fmt.Println("Could not get user info: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseStudentToStudent(student))
}

func HandlerUpdateTutorProfile(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	tutor, ok := r.Context().Value(contextKeys.TutorKey).(domain.Tutor)
	if !ok {
		fmt.Println("Tutor profile cannot be found in context.")
		util.RespondWithInternalServerError(w)
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
		fmt.Println("Couldn't decode parameters", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.Username == "" || params.Name == "" || params.Email == "" {
		fmt.Println("Missing one more more required parameters.")
		util.RespondWithMissingParametersError(w)
		return
	}

	usernameTaken, err := apiCfg.DB.CheckUsernameTaken(r.Context(), params.Username)
	if err != nil {
		fmt.Println("Couldn't check if username taken: ", err)
		util.RespondWithInternalServerError(w)
		return
	} else if tutor.Username != params.Username && usernameTaken == 1 {
		fmt.Println("Username submitted clashes with existing username")
		util.RespondWithError(w, http.StatusConflict, "Username already taken")
		return
	}

	emailTaken, err := apiCfg.DB.CheckEmailTaken(r.Context(), params.Email)
	if err != nil {
		fmt.Println("Couldn't check if email taken: ", err)
		util.RespondWithInternalServerError(w)
		return
	} else if tutor.Email != params.Email && emailTaken == 1 {
		fmt.Println("Username submitted clashes with existing username")
		util.RespondWithError(w, http.StatusConflict, "Email already taken")
		return
	}

	tx, err := apiCfg.DBConn.BeginTx(r.Context(), nil)

	if err != nil {
		fmt.Println("Couldn't start transaction: ", err)
		util.RespondWithInternalServerError(w)
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
		fmt.Println("Couldn't update tutor profile", err)
		util.RespondWithInternalServerError(w)
		return
	}

	err = queries.UpdateUserProfile(r.Context(), database.UpdateUserProfileParams{
		Username: params.Username,
		UserID:   tutor.TutorID,
		Email:    tutor.Email,
	})
	if err != nil {
		fmt.Println("Couldn't update user profile", err)
		util.RespondWithInternalServerError(w)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Couldn't commit transaction: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	tutorSubjects, err := apiCfg.DB.GetTutorSubjects(r.Context(), tutor.TutorID)
	if err != nil {
		fmt.Println("Couldn't get tutor subjects: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseTutorToTutor(updatedTutor, tutorSubjects))
}

func HandlerUpdateTutorPassword(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
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
		fmt.Println("Tutor profile not found in context.")
		util.RespondWithInternalServerError(w)
		return
	}
	util.RespondWithJSON(w, http.StatusOK, tutor)
}

func HandlerGetAvailableQuestions(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	tutor, ok := r.Context().Value(contextKeys.TutorKey).(domain.Tutor)
	if !ok {
		fmt.Println("Tutor profile not found in context.")
		util.RespondWithInternalServerError(w)
		return
	}

	databaseChats, err := apiCfg.DB.TutorGetAvailableQuestions(r.Context(), tutor.TutorID)
	if err != nil {
		fmt.Println("Could not retrieve any of the available questions: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	var chats []domain.Chat
	for _, chat := range databaseChats {
		chats = append(chats, domain.DatabaseChatToChat(chat, []int32{}))
	}

	util.RespondWithJSON(w, http.StatusOK, chats)
}

func HandlerAcceptQuestion(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	tutor, ok := r.Context().Value(contextKeys.TutorKey).(domain.Tutor)
	if !ok {
		fmt.Println("Tutor profile not found in context.")
		util.RespondWithInternalServerError(w)
		return
	}

	type parameters struct {
		ChatID int32 `json:"chat_id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Couldn't decode parameters", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.ChatID == 0 {
		fmt.Println("Missing ChatID parameter.")
		util.RespondWithMissingParametersError(w)
		return
	}

	err = apiCfg.DB.TutorAcceptQuestion(r.Context(), database.TutorAcceptQuestionParams{
		TutorID: uuid.NullUUID{
			UUID:  tutor.TutorID,
			Valid: true,
		},
		ChatID: params.ChatID})
	if err != nil {
		fmt.Println("Could not assign chat to tutor: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, "Question accepted.")
}

func HandlerUpdateChatTopics(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	chatIDString := chi.URLParam(r, "chatID")
	chatID, err := strconv.ParseInt(chatIDString, 10, 32)
	if err != nil {
		fmt.Println("Invalid chat ID: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	type parameters struct {
		Topics []int32 `json:"topics"`
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Couldn't decode parameters", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if len(params.Topics) == 0 {
		fmt.Println("No topics were selected (Zero topics passed as parameter).")
		util.RespondWithMissingParametersError(w)
		return
	}

	tx, err := apiCfg.DBConn.BeginTx(r.Context(), nil)
	if err != nil {
		fmt.Println("Couldn't start transaction: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	defer tx.Rollback()

	queries := apiCfg.DB.WithTx(tx)

	for topicID := range params.Topics {
		err = queries.AddChatTopic(r.Context(), database.AddChatTopicParams{
			ChatID:  int32(chatID),
			TopicID: int32(topicID)},
		)
		if err != nil {
			fmt.Println("Couldn't create tutor-subject relationship", err)
			util.RespondWithInternalServerError(w)
			return
		}
	}

	chat, err := apiCfg.DB.GetChatById(r.Context(), int32(chatID))
	if err != nil {
		fmt.Println("Could not retrieve updated chat: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Couldn't commit transaction: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseChatToChat(chat, params.Topics))
}

func HandlerConfigNewChat(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	tutor, ok := r.Context().Value(contextKeys.TutorKey).(domain.Tutor)
	if !ok {
		fmt.Println("Tutor profile not found in context.")
		util.RespondWithInternalServerError(w)
		return
	}

	chatIDString := chi.URLParam(r, "chatID")
	chatID, err := strconv.ParseInt(chatIDString, 10, 32)
	if err != nil {
		fmt.Println("Invalid chat ID: ", err)
		util.RespondWithInternalServerError(w)
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
