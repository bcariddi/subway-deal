package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8192
)

// Client represents a WebSocket connection
type Client struct {
	ID         string
	PlayerName string
	RoomID     string
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	handler    MessageHandler
}

// MessageHandler processes incoming messages
type MessageHandler interface {
	HandleMessage(client *Client, messageType string, data json.RawMessage)
	HandleDisconnect(client *Client)
}

// Message represents a WebSocket message
type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func NewClient(id string, hub *Hub, conn *websocket.Conn, handler MessageHandler) *Client {
	return &Client{
		ID:      id,
		hub:     hub,
		conn:    conn,
		send:    make(chan []byte, 256),
		handler: handler,
	}
}

// ReadPump pumps messages from the WebSocket connection to the handler
func (c *Client) ReadPump() {
	defer func() {
		c.handler.HandleDisconnect(c)
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		c.handler.HandleMessage(c, msg.Type, msg.Data)
	}
}

// WritePump pumps messages from the send channel to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendJSON sends a JSON message to the client
func (c *Client) SendJSON(msgType string, data interface{}) error {
	msg := map[string]interface{}{
		"type": msgType,
		"data": data,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case c.send <- bytes:
	default:
		return nil // Buffer full, drop message
	}
	return nil
}
