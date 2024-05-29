package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/domain"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

func (apiCfg *ApiConfig) HandlerTutorRegistration(w http.ResponseWriter, r *http.Request) {
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

func (apiCfg *ApiConfig) HandlerGetStudent(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		StudentID string `json:"id"`
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

func (apiCfg *ApiConfig) HandlerUpdateTutorProfile(w http.ResponseWriter, r *http.Request, tutor database.Tutor) {
	type parameters struct {
		Username string `json: "username"`
		Name     string `json: "name"`
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

	tutor, err = apiCfg.DB.UpdateTutorProfile(r.Context(), database.UpdateTutorProfileParams{
		TutorID:  tutor.TutorID,
		Username: params.Username,
		Name:     params.Name,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update tutor profile")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseTutorToTutor(tutor))
}

func (apiCfg *ApiConfig) HandlerUpdateTutorPassword(w http.ResponseWriter, r *http.Request, tutor database.Tutor) {
	type parameters struct {
		OldPassword string `json: "old_password"`
		NewPassword string `json: "new_password"`
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

	passwordMatched := util.CheckPasswordHash(hashedOldPassword, tutor.HashedPassword)
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

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseTutorToTutor(tutor))
}

func (apiCfg *ApiConfig) HandlerGetTutorProfile(w http.ResponseWriter, r *http.Request, tutor database.Tutor) {
	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseTutorToTutor(tutor))
}
