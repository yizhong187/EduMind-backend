package ws

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

func (h *Handler) CreateAndJoinRoom(w http.ResponseWriter, r *http.Request, user database.User) {
	roomID := chi.URLParam(r, "roomId")

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

}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	h.hub.Rooms[params.ID] = &Room{
		ID:      params.ID,
		Name:    params.Name,
		Clients: make(map[uuid.UUID]*Client),
	}

	util.RespondWithJSON(w, http.StatusOK, params)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	roomID := chi.URLParam(r, "roomId")
	clientID := r.URL.Query().Get("userId")
	username := r.URL.Query().Get("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(w http.ResponseWriter, r *http.Request) {
	var clients []ClientRes
	roomId := chi.URLParam(r, "roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
		return
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}
