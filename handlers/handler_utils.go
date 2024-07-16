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

type Topic struct {
	SubjectID int    `json:"subject_id"`
	TopicID   int    `json:"topic_id"`
	Name      string `json:"name"`
}

type Subject struct {
	SubjectID int    `json:"subject_id"`
	Name      string `json:"name"`
}

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

	var subjectList []Subject
	for _, s := range subjectIDPairs {
		subjectList = append(subjectList, Subject{
			SubjectID: int(s.SubjectID),
			Name:      s.Name,
		})
	}

	util.RespondWithJSON(w, http.StatusOK, subjectList)
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

	var topicList []Topic
	for _, s := range topics {
		topicList = append(topicList, Topic{
			SubjectID: int(s.SubjectID),
			TopicID:   int(s.TopicID),
			Name:      s.Name,
		})
	}

	util.RespondWithJSON(w, http.StatusOK, topicList)
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

	var topicList []Topic
	for _, s := range topics {
		topicList = append(topicList, Topic{
			SubjectID: int(s.SubjectID),
			TopicID:   int(s.TopicID),
			Name:      s.Name,
		})
	}

	util.RespondWithJSON(w, http.StatusOK, topicList)
}
