package ws

import (
	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
)

type Room struct {
	ID      string                `json:"id"`
	Name    string                `json:"name"`
	Clients map[uuid.UUID]*Client `json:"clients"`
}

type Hub struct {
	DB         *database.Queries
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub(db *database.Queries) *Hub {
	return &Hub{
		DB:         db,
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	// Start an infinite loop to continuously process incoming events.
	for {
		select {
		// When client registers, check if room exists and then add client to the room.
		case cl := <-h.Register:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
				}
			}

		// When client unregisters, remove the client from the room, and close its message channel
		// If no clients left in the room, delete the room.
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
					if len(h.Rooms[cl.RoomID].Clients) == 0 {
						delete(h.Rooms, cl.RoomID)
					}
				}
			}

		// When a message is broadcasted, it is sent to all clients in the room.
		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomID]; ok {
				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
