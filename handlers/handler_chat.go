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

func (apiCfg *ApiConfig) HandlerStartNewChat(w http.ResponseWriter, r *http.Request) {
	// local struct to hold expected data from the request body
	type parameters struct {
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

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
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

	util.RespondWithJSON(w, http.StatusCreated, domain.DatabaseUserToUser(user))
}
