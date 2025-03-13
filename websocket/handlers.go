package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

var (
	clients   = make(map[*Client]bool)
	clientsMu sync.RWMutex
	broadcast = make(chan []byte)
)

func init() {
	go handleBroadcasts()
}

func handleBroadcasts() {
	for message := range broadcast {
		clientsMu.RLock()
		for client := range clients {
			go func(c *Client, msg []byte) {
				c.Mu.Lock()
				defer c.Mu.Unlock()
				if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Printf("Error writing to client: %v", err)
					clientsMu.Lock()
					delete(clients, c)
					clientsMu.Unlock()
					c.Conn.Close()
				}
			}(client, message)
		}
		clientsMu.RUnlock()
	}
}

func HandleWebSocket(c *websocket.Conn) {
	client := &Client{Conn: c}

	// Register client
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, client)
		clientsMu.Unlock()
		c.Close()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		// Echo the message back for now
		// In production, you might want to validate/process the message
		if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}

// BroadcastResourceUpdate sends a resource update to all connected clients
func BroadcastResourceUpdate(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling broadcast data: %v", err)
		return
	}

	broadcast <- jsonData
}
