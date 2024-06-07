package ws

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/yizhong187/EduMind-backend/internal/database"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       uuid.UUID `json:"id"`
	RoomID   string    `json:"roomId"`
	Username string    `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

// Continuously reads messages from the client's Message channel and
// sends them over the WebSocket connection as JSON.
func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

// Continuously reads messages from the WebSocket connection and sends them
// to the Broadcast channel of the hub, and unregisters the client if an error
// occurs or the connection is closed.
func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		// Create a background context
		ctx := context.Background()
		// Create a context with a timeout (e.g., 1 minute)
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel() // Ensure the context is cancelled to free up resources

		intRoomID, err := strconv.Atoi(c.RoomID)
		if err != nil {
			fmt.Println("Invalid room ID:", err)
			return
		}

		err = hub.DB.CreateNewMessage(ctx, database.CreateNewMessageParams{
			MessageID: uuid.New(),
			ChatID:    int32(intRoomID),
			UserID:    c.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Content:   msg.Content,
		})
		if err != nil {
			log.Printf("Failed to add new message to database: %v", err)
		}

		hub.Broadcast <- msg
	}
}
