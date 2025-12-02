package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// WebSocketClient represents a connected client
type WebSocketClient struct {
	conn *websocket.Conn
	send chan []byte
}

// WebSocketHub manages all connected clients
type WebSocketHub struct {
	clients    map[*WebSocketClient]bool
	broadcast  chan []byte
	register   chan *WebSocketClient
	unregister chan *WebSocketClient
	mu         sync.RWMutex
}

var hub = &WebSocketHub{
	clients:    make(map[*WebSocketClient]bool),
	broadcast:  make(chan []byte, 256),
	register:   make(chan *WebSocketClient),
	unregister: make(chan *WebSocketClient),
}

// Run starts the WebSocket hub
func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastSignals sends signals to all connected clients
func BroadcastSignals(signals interface{}) {
	data, err := json.Marshal(signals)
	if err != nil {
		log.Printf("Error marshaling signals: %v", err)
		return
	}
	hub.broadcast <- data
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(c *websocket.Conn) {
	client := &WebSocketClient{
		conn: c,
		send: make(chan []byte, 256),
	}

	hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	client.readPump()
}

// readPump reads messages from the WebSocket connection
func (c *WebSocketClient) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// writePump writes messages to the WebSocket connection
func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// StartSignalBroadcaster starts broadcasting signals periodically
func StartSignalBroadcaster() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Generate demo signals for testing
			signals := generateDemoSignals()
			BroadcastSignals(signals)
		}
	}()
}

// generateDemoSignals generates demo signals for testing
func generateDemoSignals() []map[string]interface{} {
	signals := []map[string]interface{}{
		{
			"type":      "BUY",
			"entry":     91234.56,
			"stopLoss":  90800.00,
			"targets": []map[string]interface{}{
				{"price": 92000.00, "rr": 1.76},
			},
			"strength":  78.5,
			"timeframe": "15m",
		},
		{
			"type":      "SELL",
			"entry":     91150.00,
			"stopLoss":  91600.00,
			"targets": []map[string]interface{}{
				{"price": 90250.00, "rr": 2.0},
			},
			"strength":  82.3,
			"timeframe": "1h",
		},
	}

	return signals
}

// WebSocketUpgrade middleware for WebSocket upgrade
func WebSocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
