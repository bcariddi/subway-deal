package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/bcariddi/subway-deal/game"
)

// GameSession holds a game instance
type GameSession struct {
	Engine      *game.Engine
	PlayerNames []string
}

// Server manages game sessions
type Server struct {
	sessions map[string]*GameSession
	mu       sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		sessions: make(map[string]*GameSession),
	}
}

func main() {
	server := NewServer()

	// API routes
	http.HandleFunc("/api/game/new", server.handleNewGame)
	http.HandleFunc("/api/game/state", server.handleGetState)
	http.HandleFunc("/api/game/action", server.handleAction)

	// Serve static files from web/ directory
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	port := "8080"
	log.Printf("Subway Deal web server starting at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// POST /api/game/new - Create new game
func (s *Server) handleNewGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerNames []string `json:"playerNames"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if len(req.PlayerNames) < 2 || len(req.PlayerNames) > 5 {
		http.Error(w, "Need 2-5 players", http.StatusBadRequest)
		return
	}

	// Create game engine
	engine := game.NewEngine("local-game", req.PlayerNames)
	engine.StartTurn()

	// Store session (single session for simplicity)
	s.mu.Lock()
	s.sessions["default"] = &GameSession{
		Engine:      engine,
		PlayerNames: req.PlayerNames,
	}
	s.mu.Unlock()

	log.Printf("New game created with players: %v", req.PlayerNames)

	// Return initial state
	s.writeGameState(w, engine)
}

// GET /api/game/state - Get current game state
func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.RLock()
	session := s.sessions["default"]
	s.mu.RUnlock()

	if session == nil {
		http.Error(w, "No active game", http.StatusNotFound)
		return
	}

	s.writeGameState(w, session.Engine)
}

// POST /api/game/action - Execute game action
func (s *Server) handleAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	session := s.sessions["default"]
	s.mu.Unlock()

	if session == nil {
		http.Error(w, "No active game", http.StatusNotFound)
		return
	}

	var req struct {
		Type     string            `json:"type"`
		PlayerID string            `json:"playerId"`
		Data     map[string]string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	engine := session.Engine

	// Handle END_TURN specially
	if req.Type == "END_TURN" {
		engine.EndTurn()
		log.Printf("Turn ended, now %s's turn", engine.State.GetCurrentPlayer().Name)
		s.writeGameState(w, engine)
		return
	}

	// Execute action
	action := game.NewAction(game.ActionType(req.Type), req.PlayerID, req.Data)
	result, err := engine.ExecuteAction(action)

	if err != nil {
		log.Printf("Action failed: %s", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Action executed: %s by %s - %s", req.Type, req.PlayerID, result.Message)

	// Return updated state with action result
	state := s.buildFullState(engine)
	state["actionResult"] = result

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (s *Server) writeGameState(w http.ResponseWriter, engine *game.Engine) {
	state := s.buildFullState(engine)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (s *Server) buildFullState(engine *game.Engine) map[string]interface{} {
	state := engine.GetPublicState()

	// Add full hand details for each player (hot-seat mode shows all)
	players := make([]map[string]interface{}, len(engine.State.Players))
	for i, p := range engine.State.Players {
		players[i] = map[string]interface{}{
			"id":           p.ID,
			"name":         p.Name,
			"hand":         serializeCards(p.Hand),
			"bank":         serializeCards(p.Bank),
			"bankTotal":    p.GetTotalMoney(),
			"properties":   serializeProperties(p),
			"completeSets": p.GetCompleteSetCount(),
		}
	}
	state["players"] = players

	// Add pending action details if in response phase
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

	// Add game over info
	if engine.State.IsGameOver() && engine.State.Winner != nil {
		state["winner"] = map[string]interface{}{
			"id":   engine.State.Winner.ID,
			"name": engine.State.Winner.Name,
		}
	}

	return state
}

func serializeCards(cards []game.Card) []map[string]interface{} {
	result := make([]map[string]interface{}, len(cards))
	for i, card := range cards {
		cardMap := map[string]interface{}{
			"id":    card.GetID(),
			"type":  card.GetType(),
			"name":  card.GetName(),
			"value": card.GetValue(),
		}

		// Add type-specific fields
		switch c := card.(type) {
		case *game.PropertyCard:
			cardMap["color"] = c.Color
			cardMap["colorHex"] = c.ColorHex
			cardMap["rent"] = c.Rent
		case *game.WildcardCard:
			cardMap["colors"] = c.Colors
			cardMap["currentColor"] = c.CurrentColor
		case *game.ActionCard:
			cardMap["effect"] = c.Effect
		case *game.RentCard:
			cardMap["colors"] = c.Colors
			cardMap["colorNames"] = c.ColorNames
			cardMap["isWildRent"] = c.IsWildRent()
		}

		result[i] = cardMap
	}
	return result
}

func serializeProperties(player *game.Player) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for color, set := range player.Properties {
		if len(set.Cards) == 0 {
			continue
		}
		result = append(result, map[string]interface{}{
			"color":        color,
			"cards":        serializeCards(cardsToSlice(set.Cards)),
			"setSize":      set.SetSize(),
			"isComplete":   set.IsComplete(),
			"rent":         set.GetRent(),
			"improvements": set.Improvements,
		})
	}
	return result
}

// Helper to convert []Card to []game.Card for serialization
func cardsToSlice(cards []game.Card) []game.Card {
	return cards
}
