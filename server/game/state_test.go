package game

import (
	"testing"
)

func TestNewGameState(t *testing.T) {
	gs := NewGameState("game-123")

	if gs.ID != "game-123" {
		t.Errorf("Expected ID game-123, got %s", gs.ID)
	}
	if gs.Phase != PhaseSetup {
		t.Errorf("Expected PhaseSetup, got %s", gs.Phase)
	}
	if gs.TurnPhase != TurnPhaseDraw {
		t.Errorf("Expected TurnPhaseDraw, got %s", gs.TurnPhase)
	}
	if gs.MaxActionsPerTurn != 3 {
		t.Errorf("Expected 3 max actions, got %d", gs.MaxActionsPerTurn)
	}
	if len(gs.Players) != 0 {
		t.Errorf("Expected 0 players, got %d", len(gs.Players))
	}
	if len(gs.Deck) != 0 {
		t.Errorf("Expected empty deck, got %d cards", len(gs.Deck))
	}
}

func TestGameStatePlayerManagement(t *testing.T) {
	gs := NewGameState("game-123")

	p1 := NewPlayer("p1", "Alice")
	p2 := NewPlayer("p2", "Bob")
	gs.Players = append(gs.Players, p1, p2)

	// GetCurrentPlayer
	current := gs.GetCurrentPlayer()
	if current != p1 {
		t.Error("First player should be current")
	}

	// GetPlayer
	found := gs.GetPlayer("p2")
	if found != p2 {
		t.Error("Failed to find player by ID")
	}

	notFound := gs.GetPlayer("nonexistent")
	if notFound != nil {
		t.Error("Should return nil for non-existent player")
	}

	// NextPlayer
	gs.NextPlayer()
	if gs.CurrentPlayerIndex != 1 {
		t.Errorf("Expected player index 1, got %d", gs.CurrentPlayerIndex)
	}
	if gs.GetCurrentPlayer() != p2 {
		t.Error("Current player should be Bob")
	}

	// Wrap around
	gs.NextPlayer()
	if gs.CurrentPlayerIndex != 0 {
		t.Error("Player index should wrap around to 0")
	}
}

func TestGameStateNextPlayerResetsState(t *testing.T) {
	gs := NewGameState("game-123")
	gs.Players = append(gs.Players, NewPlayer("p1", "Alice"), NewPlayer("p2", "Bob"))

	gs.TurnPhase = TurnPhaseActions
	gs.ActionsPlayedThisTurn = 2

	gs.NextPlayer()

	if gs.TurnPhase != TurnPhaseDraw {
		t.Errorf("Expected TurnPhaseDraw after NextPlayer, got %s", gs.TurnPhase)
	}
	if gs.ActionsPlayedThisTurn != 0 {
		t.Errorf("Expected 0 actions after NextPlayer, got %d", gs.ActionsPlayedThisTurn)
	}
}

func TestGameStateIsGameOver(t *testing.T) {
	gs := NewGameState("game-123")

	if gs.IsGameOver() {
		t.Error("New game should not be over")
	}

	gs.Phase = PhasePlaying
	if gs.IsGameOver() {
		t.Error("Playing game should not be over")
	}

	gs.Phase = PhaseFinished
	if !gs.IsGameOver() {
		t.Error("Finished game should be over")
	}
}

func TestGameStateCheckWinCondition(t *testing.T) {
	gs := NewGameState("game-123")
	p1 := NewPlayer("p1", "Alice")
	gs.Players = append(gs.Players, p1)

	// No winner initially
	if gs.CheckWinCondition() {
		t.Error("Should not have winner initially")
	}
	if gs.Winner != nil {
		t.Error("Winner should be nil")
	}

	// Complete 3 sets
	for _, color := range []string{"brown", "utility", "darkblue"} {
		set := p1.GetPropertySet(color)
		for i := 0; i < set.SetSize(); i++ {
			set.AddCard(&PropertyCard{
				BaseCard: BaseCard{ID: color + string(rune('1'+i))},
				Color:    color,
			})
		}
	}

	if !gs.CheckWinCondition() {
		t.Error("Should have winner with 3 complete sets")
	}
	if gs.Winner != p1 {
		t.Error("Winner should be Alice")
	}
	if gs.Phase != PhaseFinished {
		t.Error("Game phase should be finished")
	}
}

func TestGetCurrentPlayerEmpty(t *testing.T) {
	gs := NewGameState("game-123")

	current := gs.GetCurrentPlayer()
	if current != nil {
		t.Error("Should return nil when no players")
	}
}
