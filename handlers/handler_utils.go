package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yizhong187/EduMind-backend/contextKeys"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

func HandlerGetAllSubjects(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	subjectIDPairs, err := apiCfg.DB.GetAllSubjects(r.Context())
	if err != nil {
		fmt.Println("Could not retrieve subjects from subjects table: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, subjectIDPairs)
}

func HandlerGetSubjectTopics(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	subjectID := chi.URLParam(r, "subjectID")
	parsedSubjectID, err := strconv.ParseInt(subjectID, 10, 32)
	if err != nil {
		fmt.Println("Invalid subject ID: ", err)
		util.RespondWithError(w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	topics, err := apiCfg.DB.GetAllTopicsBySubject(r.Context(), int32(parsedSubjectID))
	if err != nil {
		fmt.Println("Could not retrieve topics from subject_topics table: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, topics)
}

func HandlerGetAllTopics(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	topics, err := apiCfg.DB.GetAllTopics(r.Context())
	if err != nil {
		fmt.Println("Could not retrieve topics from subject_topics table: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, topics)
}
