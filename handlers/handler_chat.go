package handlers

import (
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

func HandlerGetAllChats(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	user, ok := r.Context().Value(contextKeys.UserKey).(domain.User)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "User not found")
		return
	}

	databaseChats, err := apiCfg.DB.GetAllChats(r.Context(), user.ID)
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

func HandlerGetAllMessages(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	chatIDString := chi.URLParam(r, "chatID")
	chatID, err := strconv.ParseInt(chatIDString, 10, 32)

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Invalid chat ID")
		return
	}

	databaseMessages, err := apiCfg.DB.GetAllMessages(r.Context(), int32(chatID))
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve messages")
		return
	}

	var messages []domain.Message
	for _, message := range databaseMessages {
		messages = append(messages, domain.DatabaseMessageToMessage(message))
	}

	util.RespondWithJSON(w, http.StatusOK, messages)
}

func HandlerNewMessage(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	user, ok := r.Context().Value(contextKeys.UserKey).(domain.User)
	if !ok {
		util.RespondWithError(w, http.StatusInternalServerError, "User not found")
		return
	}

	chatIDString := chi.URLParam(r, "chatID")
	chatID, err := strconv.ParseInt(chatIDString, 10, 32)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Invalid chat ID")
		return
	}

	type parameters struct {
		Content string `json:"content"`
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	if params.Content == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Content is required")
		return
	}

	err = apiCfg.DB.CreateNewMessage(r.Context(), database.CreateNewMessageParams{
		MessageID: uuid.New(),
		ChatID:    int32(chatID),
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Content:   params.Content,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create new message")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, struct{}{})
}
