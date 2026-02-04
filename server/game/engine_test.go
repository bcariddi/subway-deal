package game

import "testing"

func TestNewEngine(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})

	if engine.State.ID != "test-game" {
		t.Errorf("Expected game ID test-game, got %s", engine.State.ID)
	}

	if len(engine.State.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(engine.State.Players))
	}

	if engine.State.Players[0].Name != "Alice" {
		t.Errorf("Expected first player Alice, got %s", engine.State.Players[0].Name)
	}

	if engine.State.Phase != PhasePlaying {
		t.Errorf("Expected PhasePlaying, got %s", engine.State.Phase)
	}

	// Each player should have 5 cards dealt
	for _, player := range engine.State.Players {
		if len(player.Hand) != 5 {
			t.Errorf("Expected 5 cards in hand for %s, got %d", player.Name, len(player.Hand))
		}
	}

	// Deck should have 106 - (2 players * 5 cards) = 96 cards
	expectedDeckSize := 106 - (2 * 5)
	if len(engine.State.Deck) != expectedDeckSize {
		t.Errorf("Expected %d cards in deck, got %d", expectedDeckSize, len(engine.State.Deck))
	}
}

func TestEngineStartTurn(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})

	alice := engine.State.Players[0]
	initialHandSize := len(alice.Hand)
	initialDeckSize := len(engine.State.Deck)

	engine.StartTurn()

	// Should have drawn 2 cards
	if len(alice.Hand) != initialHandSize+2 {
		t.Errorf("Expected %d cards after drawing, got %d", initialHandSize+2, len(alice.Hand))
	}

	// Deck should be smaller
	if len(engine.State.Deck) != initialDeckSize-2 {
		t.Errorf("Expected %d cards in deck, got %d", initialDeckSize-2, len(engine.State.Deck))
	}

	// Turn phase should be actions
	if engine.State.TurnPhase != TurnPhaseActions {
		t.Errorf("Expected TurnPhaseActions, got %s", engine.State.TurnPhase)
	}
}

func TestEngineStartTurnEmptyHand(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})

	alice := engine.State.Players[0]
	// Clear Alice's hand
	alice.Hand = nil

	engine.StartTurn()

	// Should have drawn 5 cards when hand was empty
	if len(alice.Hand) != 5 {
		t.Errorf("Expected 5 cards with empty hand, got %d", len(alice.Hand))
	}
}

func TestEngineEndTurn(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()

	// It's Alice's turn
	if engine.State.GetCurrentPlayer().Name != "Alice" {
		t.Error("Expected Alice's turn")
	}

	engine.EndTurn()

	// Should be Bob's turn now
	if engine.State.GetCurrentPlayer().Name != "Bob" {
		t.Error("Expected Bob's turn after EndTurn")
	}

	// Bob should have drawn cards
	bob := engine.State.Players[1]
	if len(bob.Hand) <= 5 {
		t.Error("Bob should have drawn cards")
	}
}

func TestEngineHandLimit(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()

	alice := engine.State.Players[0]

	// Give Alice more than 7 cards
	for i := 0; i < 5; i++ {
		alice.AddToHand(&MoneyCard{
			BaseCard: BaseCard{ID: "extra_" + string(rune('a'+i)), Value: 1},
		})
	}

	if len(alice.Hand) <= 7 {
		t.Fatal("Alice should have more than 7 cards for this test")
	}

	engine.EndTurn()

	// After end turn, Alice should have at most 7 cards
	if len(alice.Hand) > 7 {
		t.Errorf("Hand limit not enforced: Alice has %d cards", len(alice.Hand))
	}
}

func TestEngineExecuteAction(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()

	alice := engine.State.Players[0]

	// Find a money card in Alice's hand
	var moneyCard Card
	for _, card := range alice.Hand {
		if card.GetType() == CardTypeMoney {
			moneyCard = card
			break
		}
	}

	if moneyCard == nil {
		// Add a money card for testing
		moneyCard = &MoneyCard{BaseCard: BaseCard{ID: "test_money", Type: CardTypeMoney, Value: 1}}
		alice.AddToHand(moneyCard)
	}

	action := NewAction(ActionPlayMoney, alice.ID, map[string]string{
		"cardId": moneyCard.GetID(),
	})

	result, err := engine.ExecuteAction(action)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("Expected success, got error: %s", result.Error)
	}

	// Card should be in bank now
	if len(alice.Bank) == 0 {
		t.Error("Card should be in bank")
	}
}

func TestEngineWinCondition(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()

	alice := engine.State.Players[0]

	// Give Alice 3 complete sets
	colors := []string{"brown", "utility", "darkblue"}
	for _, color := range colors {
		set := alice.GetPropertySet(color)
		setSize := set.SetSize()
		for i := 0; i < setSize; i++ {
			set.AddCard(&PropertyCard{
				BaseCard: BaseCard{ID: color + string(rune('1'+i))},
				Color:    color,
				Rent:     map[int]int{1: 1, 2: 2},
			})
		}
	}

	// Execute any action to trigger win check
	moneyCard := &MoneyCard{BaseCard: BaseCard{ID: "test_money", Type: CardTypeMoney, Value: 1}}
	alice.AddToHand(moneyCard)

	action := NewAction(ActionPlayMoney, alice.ID, map[string]string{
		"cardId": moneyCard.GetID(),
	})

	engine.ExecuteAction(action)

	if !engine.State.IsGameOver() {
		t.Error("Game should be over")
	}
	if engine.State.Winner != alice {
		t.Error("Alice should be the winner")
	}
}

func TestEngineGetPublicState(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()

	state := engine.GetPublicState()

	if state["currentPlayer"] != engine.State.Players[0].ID {
		t.Error("Current player mismatch")
	}

	players := state["players"].([]map[string]interface{})
	if len(players) != 2 {
		t.Errorf("Expected 2 players in state, got %d", len(players))
	}

	// Should not include hand contents (just count)
	for _, p := range players {
		if _, hasHand := p["hand"]; hasHand {
			t.Error("Public state should not include hand contents")
		}
		if _, hasHandCount := p["handCount"]; !hasHandCount {
			t.Error("Public state should include hand count")
		}
	}
}

func TestEngineGetPlayerView(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()

	alice := engine.State.Players[0]
	view := engine.GetPlayerView(alice.ID)

	// Should include hand for this player
	hand, hasHand := view["hand"]
	if !hasHand {
		t.Error("Player view should include hand")
	}

	handSlice := hand.([]map[string]interface{})
	if len(handSlice) != len(alice.Hand) {
		t.Errorf("Hand size mismatch: view has %d, actual has %d", len(handSlice), len(alice.Hand))
	}
}
