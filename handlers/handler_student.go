package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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
		PhotoURL string `json:"photo_url"`
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

	var photoURL sql.NullString
	if params.PhotoURL == "" {
		photoURL = sql.NullString{String: "", Valid: false}
	} else {
		photoURL = sql.NullString{String: params.PhotoURL, Valid: true}
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
		PhotoUrl:       photoURL,
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
		fmt.Println("Student profile cannot be found in context.")
		util.RespondWithInternalServerError(w)
		return
	}

	type parameters struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		PhotoURL string `json:"photo_url"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Couldn't decode parameters: ", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.Username == "" || params.Name == "" || params.Email == "" {
		fmt.Println("Missing one or more required parameters.")
		util.RespondWithMissingParametersError(w)
		return
	}

	usernameTaken, err := apiCfg.DB.CheckUsernameTaken(r.Context(), params.Username)
	if err != nil {
		fmt.Println("Couldn't check if username taken: ", err)
		util.RespondWithInternalServerError(w)
		return
	} else if student.Username != params.Username && usernameTaken == 1 {
		fmt.Println("Username submitted clashes with existing username")
		util.RespondWithError(w, http.StatusConflict, "Username already taken")
		return
	}

	emailTaken, err := apiCfg.DB.CheckEmailTaken(r.Context(), params.Email)
	if err != nil {
		fmt.Println("Couldn't check if email taken: ", err)
		util.RespondWithInternalServerError(w)
		return
	} else if student.Email != params.Email && emailTaken == 1 {
		fmt.Println("Username submitted clashes with existing username")
		util.RespondWithError(w, http.StatusConflict, "Email already taken")
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
		fmt.Println("Couldn't start transaction: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	defer tx.Rollback()

	queries := apiCfg.DB.WithTx(tx)

	updatedStudent, err := queries.UpdateStudentProfile(r.Context(), database.UpdateStudentProfileParams{
		StudentID: student.StudentID,
		Username:  params.Username,
		Name:      params.Name,
		Email:     params.Email,
		PhotoUrl:  photoURL,
	})

	if err != nil {
		fmt.Println("Couldn't update student profile", err)
		util.RespondWithInternalServerError(w)
		return
	}

	err = queries.UpdateUserProfile(r.Context(), database.UpdateUserProfileParams{
		Username: params.Username,
		Email:    params.Email,
		UserID:   student.StudentID,
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
		fmt.Println("Student profile not found ")
		util.RespondWithInternalServerError(w)
		return
	}

	type parameters struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Couldn't decode parameters", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.OldPassword == "" || params.NewPassword == "" {
		fmt.Println("Missing one more more required parameters.")
		util.RespondWithMissingParametersError(w)
		return
	}

	databaseHashedPassword, err := apiCfg.DB.GetTutorHash(r.Context(), student.Username)
	if err != nil {
		fmt.Println("Couldn't retrieve hash: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	passwordMatched := util.CheckPasswordHash(params.OldPassword, databaseHashedPassword)
	if !passwordMatched {
		fmt.Println("Incorrect password")
		util.RespondWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	hashedNewPassword, err := util.HashPassword(params.NewPassword)
	if err != nil {
		fmt.Println("Hashing password went wrong: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	err = apiCfg.DB.UpdateStudentPassword(r.Context(), database.UpdateStudentPasswordParams{
		StudentID:      student.StudentID,
		HashedPassword: hashedNewPassword,
	})
	if err != nil {
		fmt.Println("Couldn't update password: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, "Password updated successfully")
}

func HandlerGetStudentProfile(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}
	username := chi.URLParam(r, "username")

	studentProfile, err := apiCfg.DB.GetStudentByUsername(r.Context(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Student profile not found", err)
			util.RespondWithError(w, http.StatusNotFound, "Student profile not found")
			return
		}
		fmt.Println("Couldn't retrieve student profile", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseStudentToStudent(studentProfile))
}

func HandlerGetStudentProfileById(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("Couldn't decode parameters", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.StudentID == "" {
		fmt.Println("Missing student_id parameter.")
		util.RespondWithMissingParametersError(w)
		return
	}

	parsedUUID, err := uuid.Parse(params.StudentID)
	if err != nil {
		fmt.Println("Invalid UUID", err)
		util.RespondWithInternalServerError(w)
		return
	}

	studentProfile, err := apiCfg.DB.GetStudentById(r.Context(), parsedUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Student profile not found", err)
			util.RespondWithError(w, http.StatusNotFound, "Student profile not found")
			return
		}
		fmt.Println("Couldn't retrieve student profile", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseStudentToStudent(studentProfile))
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
		fmt.Println("Student profile not found in context.")
		util.RespondWithInternalServerError(w)
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
		fmt.Println("Couldn't decode parameters", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.SubjectID == 0 || params.Header == "" {
		fmt.Println("Missing one more more required parameters.")
		util.RespondWithMissingParametersError(w)
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
		fmt.Println("Couldn't start transaction: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	queries := apiCfg.DB.WithTx(tx)

	chat, err := queries.StudentCreateNewChat(r.Context(), database.StudentCreateNewChatParams{
		StudentID: student.StudentID,
		CreatedAt: time.Now().UTC(),
		SubjectID: params.SubjectID,
		Header:    params.Header,
		PhotoUrl:  photoURL,
	})
	if err != nil {
		fmt.Println("Couldn't create new chat: ", err)
		util.RespondWithInternalServerError(w)
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
		fmt.Println("Couldn't create new message", err)
		util.RespondWithInternalServerError(w)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Couldn't commit transaction: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, "Question submitted successfully.")
}
