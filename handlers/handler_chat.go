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
		fmt.Println("User not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	databaseChats, err := apiCfg.DB.GetAllChats(r.Context(), user.ID)
	if err != nil {
		fmt.Println("Couldn't retrieve chats", err)
		util.RespondWithInternalServerError(w)
		return
	}

	chats := []domain.Chat{}
	for _, chat := range databaseChats {
		chatTopics, err := apiCfg.DB.GetChatTopics(r.Context(), chat.ChatID)
		if err != nil {
			fmt.Println("Couldn't retrieve chat topics: ", err)
			util.RespondWithInternalServerError(w)
			return
		}
		chats = append(chats, domain.DatabaseChatToChat(chat, chatTopics))
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
		fmt.Println("Invalid chat ID: ", err)
		util.RespondWithBadRequest(w, "Invalid chat ID.")
		return
	}

	databaseMessages, err := apiCfg.DB.GetAllMessages(r.Context(), int32(chatID))
	if err != nil {
		fmt.Println("Couldn't retrieve messages", err)
		util.RespondWithInternalServerError(w)
		return
	}

	messages := []domain.Message{}
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
		fmt.Println("User not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	chatIDString := chi.URLParam(r, "chatID")
	chatID, err := strconv.ParseInt(chatIDString, 10, 32)
	if err != nil {
		fmt.Println("Invalid chat ID", err)
		util.RespondWithBadRequest(w, "Invalid chat ID.")
		return
	}

	type parameters struct {
		Content string `json:"content"`
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Couldn't decode parameters", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.Content == "" {
		fmt.Println("Missing content parameter in request: ", err)
		util.RespondWithMissingParametersError(w)
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
		fmt.Println("Couldn't create new message", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, "Message sent.")
}

func HandlerCompleteChat(w http.ResponseWriter, r *http.Request) {
	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		fmt.Println("ApiConfig not found.")
		util.RespondWithInternalServerError(w)
		return
	}

	chatIDString := chi.URLParam(r, "chatID")
	chatID, err := strconv.ParseInt(chatIDString, 10, 32)
	if err != nil {
		fmt.Println("Invalid chat ID", err)
		util.RespondWithBadRequest(w, "Invalid chat ID.")
		return
	}

	err = apiCfg.DB.CompleteChat(r.Context(), int32(chatID))

	if err != nil {
		fmt.Println("Could not mark chat as complete", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, "Chat marked as complete.")
}
