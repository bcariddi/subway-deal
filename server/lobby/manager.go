package lobby

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Lobby represents a game waiting room
type Lobby struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	HostID     string    `json:"hostId"`
	MaxPlayers int       `json:"maxPlayers"`
	Players    []Player  `json:"players"`
	Status     string    `json:"status"` // "waiting", "playing", "finished"
	GameID     string    `json:"gameId,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Player represents a player in a lobby
type Player struct {
	ClientID  string `json:"clientId"`
	Name      string `json:"name"`
	Ready     bool   `json:"ready"`
	Connected bool   `json:"connected"`
}

// Manager handles lobby operations
type Manager struct {
	lobbies map[string]*Lobby
	mu      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		lobbies: make(map[string]*Lobby),
	}
}

// CreateLobby creates a new lobby
func (m *Manager) CreateLobby(name string, hostID string, hostName string, maxPlayers int) (*Lobby, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if maxPlayers < 2 || maxPlayers > 5 {
		maxPlayers = 5
	}

	lobby := &Lobby{
		ID:         uuid.New().String(),
		Name:       name,
		HostID:     hostID,
		MaxPlayers: maxPlayers,
		Players: []Player{
			{ClientID: hostID, Name: hostName, Ready: true, Connected: true},
		},
		Status:    "waiting",
		CreatedAt: time.Now(),
	}

	m.lobbies[lobby.ID] = lobby
	return lobby, nil
}

// GetLobby retrieves a lobby by ID
func (m *Manager) GetLobby(id string) (*Lobby, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lobby, ok := m.lobbies[id]
	if !ok {
		return nil, errors.New("lobby not found")
	}
	return lobby, nil
}

// GetOpenLobbies returns all lobbies that are waiting for players
func (m *Manager) GetOpenLobbies() []*Lobby {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lobbies := make([]*Lobby, 0)
	for _, lobby := range m.lobbies {
		if lobby.Status == "waiting" {
			lobbies = append(lobbies, lobby)
		}
	}
	return lobbies
}

// JoinLobby adds a player to a lobby
func (m *Manager) JoinLobby(lobbyID, clientID, playerName string) (*Lobby, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	lobby, ok := m.lobbies[lobbyID]
	if !ok {
		return nil, errors.New("lobby not found")
	}

	if lobby.Status != "waiting" {
		return nil, errors.New("lobby is not accepting players")
	}

	if len(lobby.Players) >= lobby.MaxPlayers {
		return nil, errors.New("lobby is full")
	}

	// Check if already in lobby
	for _, p := range lobby.Players {
		if p.ClientID == clientID {
			return nil, errors.New("already in lobby")
		}
	}

	lobby.Players = append(lobby.Players, Player{
		ClientID:  clientID,
		Name:      playerName,
		Ready:     false,
		Connected: true,
	})

	return lobby, nil
}

// LeaveLobby removes a player from a lobby
func (m *Manager) LeaveLobby(lobbyID, clientID string) (*Lobby, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	lobby, ok := m.lobbies[lobbyID]
	if !ok {
		return nil, errors.New("lobby not found")
	}

	// Find and remove player
	for i, p := range lobby.Players {
		if p.ClientID == clientID {
			lobby.Players = append(lobby.Players[:i], lobby.Players[i+1:]...)
			break
		}
	}

	// If host left, assign new host
	if lobby.HostID == clientID && len(lobby.Players) > 0 {
		lobby.HostID = lobby.Players[0].ClientID
	}

	// Delete empty lobbies
	if len(lobby.Players) == 0 {
		delete(m.lobbies, lobbyID)
		return nil, nil
	}

	return lobby, nil
}

// SetPlayerReady updates a player's ready status
func (m *Manager) SetPlayerReady(lobbyID, clientID string, ready bool) (*Lobby, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	lobby, ok := m.lobbies[lobbyID]
	if !ok {
		return nil, errors.New("lobby not found")
	}

	for i, p := range lobby.Players {
		if p.ClientID == clientID {
			lobby.Players[i].Ready = ready
			break
		}
	}

	return lobby, nil
}

// CanStart checks if lobby can start a game
func (m *Manager) CanStart(lobbyID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lobby, ok := m.lobbies[lobbyID]
	if !ok {
		return false
	}

	if lobby.Status != "waiting" || len(lobby.Players) < 2 {
		return false
	}

	for _, p := range lobby.Players {
		if !p.Ready {
			return false
		}
	}

	return true
}

// StartGame marks lobby as playing
func (m *Manager) StartGame(lobbyID, gameID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	lobby, ok := m.lobbies[lobbyID]
	if !ok {
		return errors.New("lobby not found")
	}

	lobby.Status = "playing"
	lobby.GameID = gameID
	return nil
}

// FindLobbyByClient finds the lobby a client is in
func (m *Manager) FindLobbyByClient(clientID string) *Lobby {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, lobby := range m.lobbies {
		for _, p := range lobby.Players {
			if p.ClientID == clientID {
				return lobby
			}
		}
	}
	return nil
}

// DeleteLobby removes a lobby
func (m *Manager) DeleteLobby(lobbyID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.lobbies, lobbyID)
}
