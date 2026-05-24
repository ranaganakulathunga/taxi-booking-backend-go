package realtime

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RideUpdate represents a real-time ride status update
type RideUpdate struct {
	RideID     uint    `json:"ride_id"`
	Status     string  `json:"status"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	DriverName string  `json:"driver_name,omitempty"`
	Message    string  `json:"message,omitempty"`
}

// Client represents a WebSocket client
type Client struct {
	ID   string
	Conn *websocket.Conn
	Hub  *Hub
	Send chan *RideUpdate
}

// Hub maintains active client connections and broadcasts messages
type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan *RideUpdate
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *RideUpdate),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

// Run starts the hub's event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mu.Lock()
			h.Clients[client] = true
			h.Mu.Unlock()
			log.Printf("Client %s registered. Total clients: %d", client.ID, len(h.Clients))

		case client := <-h.Unregister:
			h.Mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				h.Mu.Unlock()
				log.Printf("Client %s unregistered. Total clients: %d", client.ID, len(h.Clients))
			} else {
				h.Mu.Unlock()
			}

		case message := <-h.Broadcast:
			h.Mu.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					// Client's send channel is full, skip
					go func(c *Client) {
						h.Unregister <- c
					}(client)
				}
			}
			h.Mu.RUnlock()
		}
	}
}

// BroadcastUpdate broadcasts a ride update to all connected clients
func (h *Hub) BroadcastUpdate(update *RideUpdate) {
	h.Broadcast <- update
}

// ClientRead reads messages from the client connection
func (c *Client) ClientRead() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Time{})
	for {
		message := &RideUpdate{}
		err := c.Conn.ReadJSON(message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}
		// Optionally process incoming messages from client
		log.Printf("Message from client %s: %v", c.ID, message)
	}
}

// ClientWrite writes messages to the client connection
func (c *Client) ClientWrite() {
	defer c.Conn.Close()

	for message := range c.Send {
		err := c.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error writing to client %s: %v", c.ID, err)
			return
		}
	}
}
