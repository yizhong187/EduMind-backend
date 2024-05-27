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
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// local struct to hold expected data from the request body
	type parameters struct {
		Name     string
		Password string
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		Name:           params.Name,
		Valid:          true,
		HashedPassword: string(hashedPassword),
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseUserToUser(user))
}

func (apiCfg *ApiConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID string
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	parsedUUID, err := uuid.Parse(params.ID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Invalid UUID")
		return
	}

	user, err := apiCfg.DB.GetUser(r.Context(), parsedUUID)

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Could not get user info")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, domain.DatabaseUserToUser(user))
}
