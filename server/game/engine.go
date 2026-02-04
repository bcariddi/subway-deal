package game

import (
	"fmt"
	"math/rand"
)

// Engine orchestrates the game
type Engine struct {
	State     *GameState
	validator *Validator
	executor  *Executor
}

// NewEngine creates a new game with the given players
func NewEngine(gameID string, playerNames []string) *Engine {
	state := NewGameState(gameID)

	// Create players
	for i, name := range playerNames {
		player := NewPlayer(fmt.Sprintf("player_%d", i), name)
		state.Players = append(state.Players, player)
	}

	engine := &Engine{
		State:     state,
		validator: NewValidator(state),
		executor:  NewExecutor(state),
	}

	// Initialize deck and deal
	engine.initializeDeck()
	engine.shuffleDeck()
	engine.dealInitialHands()

	state.Phase = PhasePlaying

	return engine
}

func (e *Engine) initializeDeck() {
	e.State.Deck = CreateAllCards()
}

func (e *Engine) shuffleDeck() {
	deck := e.State.Deck
	for i := len(deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

func (e *Engine) dealInitialHands() {
	for _, player := range e.State.Players {
		for i := 0; i < 5; i++ {
			if len(e.State.Deck) > 0 {
				card := e.State.Deck[len(e.State.Deck)-1]
				e.State.Deck = e.State.Deck[:len(e.State.Deck)-1]
				player.AddToHand(card)
			}
		}
	}
}

// ExecuteAction validates and executes a player action
func (e *Engine) ExecuteAction(action *Action) (ActionResult, error) {
	// Validate
	validation := e.validator.Validate(action)
	if !validation.Valid {
		return ActionResult{Success: false, Error: validation.Error}, fmt.Errorf("%s", validation.Error)
	}

	// Execute
	result := e.executor.Execute(action)

	// Check win condition
	e.State.CheckWinCondition()

	return result, nil
}

// StartTurn begins a player's turn by drawing cards
func (e *Engine) StartTurn() {
	player := e.State.GetCurrentPlayer()
	e.State.TurnPhase = TurnPhaseDraw

	// Draw 2 cards (or 5 if hand is empty)
	drawCount := 2
	if len(player.Hand) == 0 {
		drawCount = 5
	}

	for i := 0; i < drawCount; i++ {
		if len(e.State.Deck) > 0 {
			card := e.State.Deck[len(e.State.Deck)-1]
			e.State.Deck = e.State.Deck[:len(e.State.Deck)-1]
			player.AddToHand(card)
		}
	}

	e.State.TurnPhase = TurnPhaseActions
}

// EndTurn finishes current player's turn and moves to next
func (e *Engine) EndTurn() {
	player := e.State.GetCurrentPlayer()

	// Enforce hand limit (7 cards)
	for len(player.Hand) > 7 {
		// Discard last card (in full game, player would choose)
		discarded := player.Hand[len(player.Hand)-1]
		player.Hand = player.Hand[:len(player.Hand)-1]
		e.State.DiscardPile = append(e.State.DiscardPile, discarded)
	}

	e.State.NextPlayer()
	e.StartTurn()
}

// GetPublicState returns state safe to send to all clients
func (e *Engine) GetPublicState() map[string]interface{} {
	players := make([]map[string]interface{}, len(e.State.Players))
	for i, p := range e.State.Players {
		players[i] = map[string]interface{}{
			"id":           p.ID,
			"name":         p.Name,
			"handCount":    len(p.Hand),
			"bankValue":    p.GetTotalMoney(),
			"properties":   e.getPlayerProperties(p),
			"completeSets": p.GetCompleteSetCount(),
		}
	}

	state := map[string]interface{}{
		"currentPlayer":         e.State.GetCurrentPlayer().ID,
		"phase":                 e.State.Phase,
		"turnPhase":             e.State.TurnPhase,
		"actionsPlayedThisTurn": e.State.ActionsPlayedThisTurn,
		"maxActionsPerTurn":     e.State.MaxActionsPerTurn,
		"deckSize":              len(e.State.Deck),
		"players":               players,
	}

	if e.State.Winner != nil {
		state["winner"] = e.State.Winner.ID
	}

	return state
}

// GetPlayerView returns state from a specific player's perspective (includes their hand)
func (e *Engine) GetPlayerView(playerID string) map[string]interface{} {
	state := e.GetPublicState()

	player := e.State.GetPlayer(playerID)
	if player != nil {
		hand := make([]map[string]interface{}, len(player.Hand))
		for i, card := range player.Hand {
			hand[i] = map[string]interface{}{
				"id":    card.GetID(),
				"type":  card.GetType(),
				"name":  card.GetName(),
				"value": card.GetValue(),
			}
		}
		state["hand"] = hand
	}

	return state
}

func (e *Engine) getPlayerProperties(p *Player) []map[string]interface{} {
	props := make([]map[string]interface{}, 0)
	for color, set := range p.Properties {
		cards := make([]map[string]interface{}, len(set.Cards))
		for i, card := range set.Cards {
			cards[i] = map[string]interface{}{
				"id":   card.GetID(),
				"name": card.GetName(),
			}
		}
		props = append(props, map[string]interface{}{
			"color":        color,
			"cards":        cards,
			"cardCount":    len(set.Cards),
			"setSize":      set.SetSize(),
			"complete":     set.IsComplete(),
			"improvements": set.Improvements,
			"rent":         set.GetRent(),
		})
	}
	return props
}
