package realtime

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

// WebSocketHandler handles WebSocket upgrade requests
func WebSocketHandler(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
			return
		}

		clientID := uuid.New().String()
		client := &Client{
			ID:   clientID,
			Conn: conn,
			Hub:  hub,
			Send: make(chan *RideUpdate, 256),
		}

		hub.Register <- client

		// Start client goroutines
		go client.ClientRead()
		go client.ClientWrite()

		log.Printf("New WebSocket client connected: %s", clientID)
	}
}
