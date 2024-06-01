package ws

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/yizhong187/EduMind-backend/contextKeys"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(contextKeys.UserKey).(database.User)

	roomID := chi.URLParam(r, "chatID")

	_, ok := h.hub.Rooms[roomID]
	if !ok {
		h.hub.Rooms[roomID] = &Room{
			ID:      roomID,
			Clients: make(map[uuid.UUID]*Client),
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	cl := &Client{
		Conn:    conn,
		Message: make(chan *Message, 10),
		ID:      user.UserID,
		RoomID:  roomID,
	}

	h.hub.Register <- cl

	go cl.writeMessage()
	cl.readMessage(h.hub)

	util.RespondWithJSON(w, http.StatusOK, nil)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	util.RespondWithJSON(w, http.StatusOK, rooms)
}

type ClientRes struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (h *Handler) GetClients(w http.ResponseWriter, r *http.Request) {
	var clients []ClientRes
	roomId := chi.URLParam(r, "roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		util.RespondWithJSON(w, http.StatusOK, clients)
		return
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	util.RespondWithJSON(w, http.StatusOK, clients)
}
