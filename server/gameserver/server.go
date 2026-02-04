package gameserver

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bcariddi/subway-deal/game"
	"github.com/bcariddi/subway-deal/lobby"
	"github.com/bcariddi/subway-deal/ws"
)

// Server orchestrates games and lobbies
type Server struct {
	hub          *ws.Hub
	lobbyManager *lobby.Manager
	games        map[string]*game.Engine
	clientToGame map[string]string // clientID -> gameID
	mu           sync.RWMutex
}

func NewServer(hub *ws.Hub) *Server {
	return &Server{
		hub:          hub,
		lobbyManager: lobby.NewManager(),
		games:        make(map[string]*game.Engine),
		clientToGame: make(map[string]string),
	}
}

// HandleMessage processes incoming WebSocket messages
func (s *Server) HandleMessage(client *ws.Client, msgType string, data json.RawMessage) {
	switch msgType {
	case "lobby:create":
		s.handleCreateLobby(client, data)
	case "lobby:join":
		s.handleJoinLobby(client, data)
	case "lobby:leave":
		s.handleLeaveLobby(client)
	case "lobby:list":
		s.handleListLobbies(client)
	case "lobby:ready":
		s.handleSetReady(client, data)
	case "lobby:start":
		s.handleStartGame(client)
	case "game:action":
		s.handleGameAction(client, data)
	default:
		log.Printf("Unknown message type: %s", msgType)
	}
}

// HandleDisconnect handles client disconnection
func (s *Server) HandleDisconnect(client *ws.Client) {
	log.Printf("Client disconnected: %s", client.ID)

	// Handle lobby disconnect
	if lob := s.lobbyManager.FindLobbyByClient(client.ID); lob != nil {
		if lob.Status == "waiting" {
			s.lobbyManager.LeaveLobby(lob.ID, client.ID)
			s.broadcastLobbyUpdate(lob.ID)
		}
		// For active games, mark player as disconnected but don't remove
	}
}

func (s *Server) handleCreateLobby(client *ws.Client, data json.RawMessage) {
	var req struct {
		Name       string `json:"name"`
		PlayerName string `json:"playerName"`
		MaxPlayers int    `json:"maxPlayers"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		client.SendJSON("error", map[string]string{"message": "invalid request"})
		return
	}

	client.PlayerName = req.PlayerName

	lob, err := s.lobbyManager.CreateLobby(req.Name, client.ID, req.PlayerName, req.MaxPlayers)
	if err != nil {
		client.SendJSON("error", map[string]string{"message": err.Error()})
		return
	}

	s.hub.JoinRoom("lobby:"+lob.ID, client)
	client.SendJSON("lobby:created", lob)
	log.Printf("Lobby created: %s by %s", lob.ID, req.PlayerName)
}

func (s *Server) handleJoinLobby(client *ws.Client, data json.RawMessage) {
	var req struct {
		LobbyID    string `json:"lobbyId"`
		PlayerName string `json:"playerName"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		client.SendJSON("error", map[string]string{"message": "invalid request"})
		return
	}

	client.PlayerName = req.PlayerName

	lob, err := s.lobbyManager.JoinLobby(req.LobbyID, client.ID, req.PlayerName)
	if err != nil {
		client.SendJSON("error", map[string]string{"message": err.Error()})
		return
	}

	s.hub.JoinRoom("lobby:"+lob.ID, client)
	s.broadcastLobbyUpdate(lob.ID)
	log.Printf("%s joined lobby %s", req.PlayerName, lob.ID)
}

func (s *Server) handleLeaveLobby(client *ws.Client) {
	lob := s.lobbyManager.FindLobbyByClient(client.ID)
	if lob == nil {
		return
	}

	lobbyID := lob.ID
	s.lobbyManager.LeaveLobby(lobbyID, client.ID)
	s.hub.LeaveRoom("lobby:"+lobbyID, client)
	s.broadcastLobbyUpdate(lobbyID)
}

func (s *Server) handleListLobbies(client *ws.Client) {
	lobbies := s.lobbyManager.GetOpenLobbies()
	client.SendJSON("lobby:list", lobbies)
}

func (s *Server) handleSetReady(client *ws.Client, data json.RawMessage) {
	var req struct {
		Ready bool `json:"ready"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		return
	}

	lob := s.lobbyManager.FindLobbyByClient(client.ID)
	if lob == nil {
		return
	}

	s.lobbyManager.SetPlayerReady(lob.ID, client.ID, req.Ready)
	s.broadcastLobbyUpdate(lob.ID)
}

func (s *Server) handleStartGame(client *ws.Client) {
	lob := s.lobbyManager.FindLobbyByClient(client.ID)
	if lob == nil {
		client.SendJSON("error", map[string]string{"message": "not in a lobby"})
		return
	}

	if lob.HostID != client.ID {
		client.SendJSON("error", map[string]string{"message": "only host can start"})
		return
	}

	if !s.lobbyManager.CanStart(lob.ID) {
		client.SendJSON("error", map[string]string{"message": "cannot start game"})
		return
	}

	// Create game engine
	playerNames := make([]string, len(lob.Players))
	for i, p := range lob.Players {
		playerNames[i] = p.Name
	}

	gameID := lob.ID // Use lobby ID as game ID
	engine := game.NewEngine(gameID, playerNames)

	s.mu.Lock()
	s.games[gameID] = engine
	for _, p := range lob.Players {
		s.clientToGame[p.ClientID] = gameID
	}
	s.mu.Unlock()

	// Update lobby status
	s.lobbyManager.StartGame(lob.ID, gameID)

	// Move clients to game room
	clients := s.hub.GetRoomClients("lobby:" + lob.ID)
	for _, c := range clients {
		s.hub.LeaveRoom("lobby:"+lob.ID, c)
		s.hub.JoinRoom("game:"+gameID, c)
	}

	// Start the game
	engine.StartTurn()

	// Broadcast initial game state to all players
	s.broadcastGameState(gameID)

	log.Printf("Game started: %s with %d players", gameID, len(playerNames))
}

func (s *Server) handleGameAction(client *ws.Client, data json.RawMessage) {
	s.mu.RLock()
	gameID, ok := s.clientToGame[client.ID]
	if !ok {
		s.mu.RUnlock()
		client.SendJSON("error", map[string]string{"message": "not in a game"})
		return
	}
	engine := s.games[gameID]
	s.mu.RUnlock()

	if engine == nil {
		client.SendJSON("error", map[string]string{"message": "game not found"})
		return
	}

	// Get player ID for this client
	playerID := s.getPlayerID(gameID, client.ID)
	if playerID == "" {
		client.SendJSON("error", map[string]string{"message": "player not found"})
		return
	}

	// Parse action
	var actionReq struct {
		Type string            `json:"type"`
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal(data, &actionReq); err != nil {
		client.SendJSON("error", map[string]string{"message": "invalid action"})
		return
	}

	// Handle end turn specially
	if actionReq.Type == "END_TURN" {
		engine.EndTurn()
		s.broadcastGameState(gameID)

		if engine.State.IsGameOver() {
			s.handleGameEnd(gameID)
		}
		return
	}

	// For responses (ACCEPT, PLAY_FARE_EVASION), we need to check if there's a pending action
	// and if the player is one of the targets
	if actionReq.Type == "ACCEPT" || actionReq.Type == "PLAY_FARE_EVASION" {
		// Response actions are allowed from targets, not just current player
	} else {
		// Verify it's this player's turn for non-response actions
		if engine.State.GetCurrentPlayer().ID != playerID {
			client.SendJSON("error", map[string]string{"message": "not your turn"})
			return
		}
	}

	// Execute action
	action := game.NewAction(game.ActionType(actionReq.Type), playerID, actionReq.Data)
	result, err := engine.ExecuteAction(action)

	if err != nil {
		client.SendJSON("error", map[string]string{"message": err.Error()})
		return
	}

	// Send result to acting player
	client.SendJSON("game:action:result", result)

	// Broadcast updated state
	s.broadcastGameState(gameID)

	// Check for game end
	if engine.State.IsGameOver() {
		s.handleGameEnd(gameID)
	}
}

func (s *Server) getPlayerID(gameID, clientID string) string {
	lob, _ := s.lobbyManager.GetLobby(gameID)
	if lob == nil {
		return ""
	}

	for i, p := range lob.Players {
		if p.ClientID == clientID {
			return fmt.Sprintf("player_%d", i)
		}
	}
	return ""
}

func (s *Server) broadcastLobbyUpdate(lobbyID string) {
	lob, err := s.lobbyManager.GetLobby(lobbyID)
	if err != nil {
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "lobby:updated",
		"data": lob,
	})
	s.hub.BroadcastToRoom("lobby:"+lobbyID, msg)
}

func (s *Server) broadcastGameState(gameID string) {
	s.mu.RLock()
	engine := s.games[gameID]
	s.mu.RUnlock()

	if engine == nil {
		return
	}

	lob, _ := s.lobbyManager.GetLobby(gameID)
	if lob == nil {
		return
	}

	// Send personalized state to each player
	for i, p := range lob.Players {
		playerID := fmt.Sprintf("player_%d", i)
		player := engine.State.GetPlayer(playerID)

		state := engine.GetPublicState()
		state["yourId"] = playerID
		state["yourHand"] = s.serializeHand(player.Hand)

		// Add pending action info for response phase
		if engine.State.HasPendingAction() {
			pending := engine.State.PendingAction
			state["pendingAction"] = map[string]interface{}{
				"type":         pending.Action.Type,
				"sourcePlayer": pending.SourcePlayerID,
				"targets":      engine.State.GetPendingTargets(),
				"rentAmount":   pending.RentAmount,
				"rentColor":    pending.RentColor,
			}
		}

		msg, _ := json.Marshal(map[string]interface{}{
			"type": "game:state",
			"data": state,
		})
		s.hub.SendToClient(p.ClientID, msg)
	}
}

func (s *Server) serializeHand(hand []game.Card) []map[string]interface{} {
	cards := make([]map[string]interface{}, len(hand))
	for i, card := range hand {
		cardData := map[string]interface{}{
			"id":    card.GetID(),
			"type":  card.GetType(),
			"name":  card.GetName(),
			"value": card.GetValue(),
		}

		// Add type-specific fields
		switch c := card.(type) {
		case *game.PropertyCard:
			cardData["color"] = c.Color
		case *game.WildcardCard:
			cardData["colors"] = c.Colors
			cardData["currentColor"] = c.CurrentColor
		case *game.ActionCard:
			cardData["effect"] = c.Effect
		case *game.RentCard:
			cardData["colors"] = c.Colors
			cardData["isWildRent"] = c.IsWildRent()
		}

		cards[i] = cardData
	}
	return cards
}

func (s *Server) handleGameEnd(gameID string) {
	s.mu.RLock()
	engine := s.games[gameID]
	s.mu.RUnlock()

	if engine == nil || engine.State.Winner == nil {
		return
	}

	// Broadcast game end
	msg, _ := json.Marshal(map[string]interface{}{
		"type": "game:ended",
		"data": map[string]interface{}{
			"winner": engine.State.Winner.Name,
		},
	})
	s.hub.BroadcastToRoom("game:"+gameID, msg)

	// Cleanup after delay
	go func() {
		time.Sleep(30 * time.Second)
		s.mu.Lock()
		delete(s.games, gameID)
		// Clean up client mappings
		for clientID, gID := range s.clientToGame {
			if gID == gameID {
				delete(s.clientToGame, clientID)
			}
		}
		s.mu.Unlock()
		s.lobbyManager.DeleteLobby(gameID)
	}()
}
