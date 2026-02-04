package ws

import (
	"sync"
)

// Hub maintains active connections and broadcasts messages
type Hub struct {
	// Registered connections by client ID
	clients map[string]*Client

	// Connections grouped by room (lobby or game)
	rooms map[string]map[string]*Client

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe access
	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		rooms:      make(map[string]map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Register adds a client to the hub (called from main)
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)

				// Remove from all rooms
				for roomID, room := range h.rooms {
					delete(room, client.ID)
					if len(room) == 0 {
						delete(h.rooms, roomID)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// JoinRoom adds a client to a room
func (h *Hub) JoinRoom(roomID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[string]*Client)
	}
	h.rooms[roomID][client.ID] = client
	client.RoomID = roomID
}

// LeaveRoom removes a client from a room
func (h *Hub) LeaveRoom(roomID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[roomID]; ok {
		delete(room, client.ID)
		if len(room) == 0 {
			delete(h.rooms, roomID)
		}
	}
	client.RoomID = ""
}

// BroadcastToRoom sends a message to all clients in a room
func (h *Hub) BroadcastToRoom(roomID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.rooms[roomID]; ok {
		for _, client := range room {
			select {
			case client.send <- message:
			default:
				// Client buffer full, skip
			}
		}
	}
}

// SendToClient sends a message to a specific client
func (h *Hub) SendToClient(clientID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clients[clientID]; ok {
		select {
		case client.send <- message:
		default:
			// Client buffer full
		}
	}
}

// GetRoomClients returns all clients in a room
func (h *Hub) GetRoomClients(roomID string) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := make([]*Client, 0)
	if room, ok := h.rooms[roomID]; ok {
		for _, client := range room {
			clients = append(clients, client)
		}
	}
	return clients
}
