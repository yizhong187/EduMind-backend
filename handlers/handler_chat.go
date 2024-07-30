package handlers

import (
	"database/sql"
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

func HandlerUpdateRating(w http.ResponseWriter, r *http.Request) {
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

	chat, err := apiCfg.DB.GetChatById(r.Context(), int32(chatID))
	if err != nil {
		fmt.Println("Could not retrieve chat: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	if !chat.Completed {
		fmt.Println("Chat not completed.")
		util.RespondWithError(w, http.StatusBadRequest, "Chat not completed.")
		return
	}

	if chat.Rating.Valid {
		fmt.Println("Chat has previously been rated.")
		util.RespondWithError(w, http.StatusConflict, "Chat has previously been rated.")
		return
	}

	type parameters struct {
		Rating int32 `json:"rating"`
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Couldn't decode parameters", err)
		util.RespondWithInternalServerError(w)
		return
	}
	defer r.Body.Close()

	if params.Rating < 1 || params.Rating > 5 {
		fmt.Println("Invalid rating: ", params.Rating)
		util.RespondWithError(w, http.StatusBadRequest, "Invalid rating or missing rating parameter.")
		return
	}

	tx, err := apiCfg.DBConn.BeginTx(r.Context(), nil)
	if err != nil {
		fmt.Println("Couldn't start transaction: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	defer tx.Rollback()

	queries := apiCfg.DB.WithTx(tx)

	chat, err = queries.UpdateChatRating(r.Context(), database.UpdateChatRatingParams{
		ChatID: int32(chatID),
		Rating: sql.NullInt32{
			Int32: params.Rating,
			Valid: true,
		},
	})
	if err != nil {
		fmt.Println("Could not update chat rating: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	tutor, err := queries.GetTutorById(r.Context(), chat.TutorID.UUID)
	if err != nil {
		fmt.Println("Could not retrieve tutor info: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	var newRating sql.NullFloat64
	if tutor.RatingCount == 0 {
		newRating = sql.NullFloat64{
			Float64: float64(params.Rating),
			Valid:   true,
		}
	} else {
		newRating = sql.NullFloat64{
			Float64: (tutor.Rating.Float64*float64(tutor.RatingCount) + float64(params.Rating)) / (tutor.Rating.Float64 + 1),
			Valid:   true,
		}
	}

	err = queries.UpdateTutorRatings(r.Context(), database.UpdateTutorRatingsParams{
		TutorID:     tutor.TutorID,
		Rating:      newRating,
		RatingCount: tutor.RatingCount + 1,
	})
	if err != nil {
		fmt.Println("Could not update tutor rating: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Couldn't commit transaction: ", err)
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, "Successfully rated chat and tutor.")
}
