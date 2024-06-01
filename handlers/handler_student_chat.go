package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/domain"
	"github.com/yizhong187/EduMind-backend/internal/util"
	"github.com/yizhong187/EduMind-backend/middlewares"
)

func HandlerStartNewChat(w http.ResponseWriter, r *http.Request, student database.Student) {
	apiCfg := r.Context().Value(middlewares.ConfigKey).(config.ApiConfig)

	// local struct to hold expected data from the request body
	type parameters struct {
		Subject string `json:"subject"`
		Header  string `json:"header"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	chat, err := apiCfg.DB.CreateNewChat(r.Context(), database.CreateNewChatParams{
		StudentID: student.StudentID,
		CreatedAt: time.Now().UTC(),
		Subject:   params.Subject,
		Header:    params.Header,
	})
	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create new chat")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, domain.DatabaseChatToChat(chat))
}

func HandlerStudentGetAllChats(w http.ResponseWriter, r *http.Request, student database.Student) {
	apiCfg := r.Context().Value(middlewares.ConfigKey).(config.ApiConfig)

	databaseChats, err := apiCfg.DB.StudentGetAllChats(r.Context(), student.StudentID)
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

func HandlerStudentGetAllMessages(w http.ResponseWriter, r *http.Request, student database.Student) {
	apiCfg := r.Context().Value(middlewares.ConfigKey).(config.ApiConfig)

	chatIDString := chi.URLParam(r, "chatID")
	chatID, err := strconv.ParseInt(chatIDString, 10, 32)
	if err != nil {
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

func HandlerStudentNewMessage(w http.ResponseWriter, r *http.Request, student database.Student) {
	apiCfg := r.Context().Value(middlewares.ConfigKey).(config.ApiConfig)

	chatIDString := chi.URLParam(r, "chatID")
	chatID, err := strconv.ParseInt(chatIDString, 10, 32)
	if err != nil {
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

	err = apiCfg.DB.CreateNewMessage(r.Context(), database.CreateNewMessageParams{
		MessageID: uuid.New(),
		ChatID:    int32(chatID),
		UserID:    student.StudentID,
		CreatedAt: time.Now().UTC(),
		Content:   params.Content,
	})

	if err != nil {
		fmt.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create new message")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, nil)
}
