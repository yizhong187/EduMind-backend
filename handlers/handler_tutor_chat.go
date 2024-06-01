package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/domain"
	"github.com/yizhong187/EduMind-backend/internal/util"
	"github.com/yizhong187/EduMind-backend/middlewares"
)

func HandlerTutorGetAllChats(w http.ResponseWriter, r *http.Request, tutor database.Tutor) {
	apiCfg := r.Context().Value(middlewares.ConfigKey).(config.ApiConfig)

	databaseChats, err := apiCfg.DB.TutorGetAllChats(r.Context(), uuid.NullUUID{
		UUID:  tutor.TutorID,
		Valid: tutor.TutorID != uuid.Nil,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chats")
		return
	}

	var chats []domain.Chat
	for _, chat := range databaseChats {
		chats = append(chats, domain.DatabaseChatToChat(chat))
	}

	util.RespondWithJSON(w, http.StatusOK, chats)
}

func HandlerConfigNewChat(w http.ResponseWriter, r *http.Request, tutor database.Tutor) {
	apiCfg := r.Context().Value(middlewares.ConfigKey).(config.ApiConfig)

	type parameters struct {
		Topic  string `json:"topic"`
		ChatID int32  `json:"chat_id`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	chat, err := apiCfg.DB.TutorUpdateChat(r.Context(), database.TutorUpdateChatParams{
		TutorID: uuid.NullUUID{
			UUID:  tutor.TutorID,
			Valid: tutor.TutorID != uuid.Nil,
		},
		Topic: sql.NullString{
			String: params.Topic,
			Valid:  params.Topic != "",
		},
		ChatID: params.ChatID,
	})

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't update chat topic")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, chat)
}

func HandlerTutorGetAllMessages(w http.ResponseWriter, r *http.Request, student database.Student) {
	apiCfg := r.Context().Value(middlewares.ConfigKey).(config.ApiConfig)

	type parameters struct {
		ChatID int32 `json:"chat_id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	databaseMessages, err := apiCfg.DB.GetAllMessages(r.Context(), params.ChatID)

	var messages []domain.Message
	for _, message := range databaseMessages {
		messages = append(messages, domain.DatabaseMessageToMessage(message))
	}

	util.RespondWithJSON(w, http.StatusOK, messages)
}
