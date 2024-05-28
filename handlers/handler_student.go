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

func (apiCfg *ApiConfig) HandlerStudentRegistration(w http.ResponseWriter, r *http.Request) {
	// local struct to hold expected data from the request body
	type parameters struct {
		Username string `json: "username"`
		Password string `json: "password"`
		Name     string `json: "name"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	hashedPassword, err := util.HashPassword(params.Password)
	if err != nil {
		log.Fatal(err)
	}

	student, err := apiCfg.DB.CreateNewStudent(r.Context(), database.CreateNewStudentParams{
		StudentID:      uuid.New(),
		CreatedAt:      time.Now().UTC(),
		Username:       params.Username,
		Name:           params.Name,
		Valid:          true,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, domain.DatabaseStudentToStudent(student))
}

func (apiCfg *ApiConfig) HandlerGetStudentProfile(w http.ResponseWriter, r *http.Request, student database.Student) {
	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseStudentToStudent(student))
}
