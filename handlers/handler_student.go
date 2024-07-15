package handlers

import (
	"database/sql"
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
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	// local struct to hold expected data from the request body
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Couldn't decode parameters: ", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.Username == "" || params.Password == "" || params.Email == "" || params.Name == "" {
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
		fmt.Println("Username submitted clashes with existing username")
		util.RespondWithError(w, http.StatusConflict, "Username already taken")
		return
	}

	emailTaken, err := apiCfg.DB.CheckEmailTaken(r.Context(), params.Email)
	if err != nil {
		fmt.Println("Couldn't check if email taken: ", err)
		util.RespondWithInternalServerError(w)
		return
	} else if emailTaken == 1 {
		fmt.Println("Email submitted clashes with existing email")
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

	studentUUID := uuid.New()

	err = queries.InsertNewUser(r.Context(), database.InsertNewUserParams{
		UserID:   studentUUID,
		Username: params.Username,
		Email:    params.Email,
		UserType: "student",
	})

	if err != nil {
		fmt.Println("Couldn't insert into user database: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	_, err = queries.CreateNewStudent(r.Context(), database.CreateNewStudentParams{
		StudentID:      studentUUID,
		CreatedAt:      time.Now().UTC(),
		Username:       params.Username,
		Email:          params.Email,
		Name:           params.Name,
		Valid:          true,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		fmt.Println("Couldn't create new student: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Couldn't commit transaction: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, "Registration successful.")
}

func HandlerUpdateStudentProfile(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	student, ok := r.Context().Value(contextKeys.StudentKey).(domain.Student)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Student profile not found")
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
	} else if student.Username != params.Username && usernameTaken == 1 {
		util.RespondWithError(w, http.StatusConflict, "Username taken")
		return
	}

	emailTaken, err := apiCfg.DB.CheckEmailTaken(r.Context(), params.Email)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't check if email taken")
		return
	} else if student.Email != params.Email && emailTaken == 1 {
		util.RespondWithError(w, http.StatusConflict, "Email taken")
		return
	}

	updatedStudent, err := apiCfg.DB.UpdateStudentProfile(r.Context(), database.UpdateStudentProfileParams{
		StudentID: student.StudentID,
		Username:  params.Username,
		Name:      params.Name,
		Email:     params.Email,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update student profile")
		return
	}

	err = apiCfg.DB.UpdateUserProfile(r.Context(), database.UpdateUserProfileParams{
		Username: params.Username,
		Email:    params.Email,
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
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
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
		fmt.Println("Student profile not found in context.")
		util.RespondWithInternalServerError(w)
		return
	}
	util.RespondWithJSON(w, http.StatusOK, student)
}

// TO BE DEPRECATED
func HandlerStartNewChat(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	student, ok := r.Context().Value(contextKeys.StudentKey).(domain.Student)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	// local struct to hold expected data from the request body
	type parameters struct {
		SubjectID int32  `json:"subject_id"`
		Header    string `json:"header"`
		PhotoURL  string `json:"photo_url"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	if params.SubjectID == 0 {
		util.RespondWithError(w, http.StatusBadRequest, "Subject is required")
		return
	} else if params.Header == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Header is required")
		return
	}

	var photoURL sql.NullString
	if params.PhotoURL == "" {
		photoURL = sql.NullString{String: "", Valid: false}
	} else {
		photoURL = sql.NullString{String: params.PhotoURL, Valid: true}
	}

	chat, err := apiCfg.DB.CreateNewChat(r.Context(), database.CreateNewChatParams{
		StudentID: student.StudentID,
		CreatedAt: time.Now().UTC(),
		SubjectID: params.SubjectID,
		Header:    params.Header,
		PhotoUrl:  photoURL,
	})
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create new chat")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, domain.DatabaseChatToChat(chat))
}

func HandlerNewQuestion(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	student, ok := r.Context().Value(contextKeys.StudentKey).(domain.Student)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	type parameters struct {
		SubjectID int32  `json:"subject_id"`
		Header    string `json:"header"`
		PhotoURL  string `json:"photo_url"`
		Content   string `json:"content"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	if params.SubjectID == 0 {
		util.RespondWithError(w, http.StatusBadRequest, "Subject is required")
		return
	} else if params.Header == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Header is required")
		return
	}

	var photoURL sql.NullString
	if params.PhotoURL == "" {
		photoURL = sql.NullString{String: "", Valid: false}
	} else {
		photoURL = sql.NullString{String: params.PhotoURL, Valid: true}
	}

	tx, err := apiCfg.DBConn.BeginTx(r.Context(), nil)

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't start transaction")
		return
	}
	defer tx.Rollback()

	queries := apiCfg.DB.WithTx(tx)

	chat, err := queries.CreateNewChat(r.Context(), database.CreateNewChatParams{
		StudentID: student.StudentID,
		CreatedAt: time.Now().UTC(),
		SubjectID: params.SubjectID,
		Header:    params.Header,
		PhotoUrl:  photoURL,
	})
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create new chat")
		return
	}

	err = queries.CreateNewMessage(r.Context(), database.CreateNewMessageParams{
		MessageID: uuid.New(),
		ChatID:    int32(chat.ChatID),
		UserID:    student.StudentID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Content:   params.Content,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create new message")
		return
	}

	if err := tx.Commit(); err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't commit transaction")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, struct{}{})
}
