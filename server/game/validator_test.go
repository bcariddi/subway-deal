package game

import "testing"

func setupTestGame() *Engine {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	return engine
}

func TestValidatorDrawCards(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	// Don't start turn yet - should be in draw phase
	validator := NewValidator(engine.State)

	alice := engine.State.Players[0]

	// Valid draw
	action := NewAction(ActionDrawCards, alice.ID, nil)
	result := validator.Validate(action)
	if !result.Valid {
		t.Errorf("Expected valid, got: %s", result.Error)
	}

	// Wrong player
	bob := engine.State.Players[1]
	action = NewAction(ActionDrawCards, bob.ID, nil)
	result = validator.Validate(action)
	if result.Valid {
		t.Error("Should not allow Bob to draw on Alice's turn")
	}

	// After starting turn, can't draw again
	engine.StartTurn()
	action = NewAction(ActionDrawCards, alice.ID, nil)
	result = validator.Validate(action)
	if result.Valid {
		t.Error("Should not allow drawing during action phase")
	}
}

func TestValidatorPlayProperty(t *testing.T) {
	engine := setupTestGame()
	validator := NewValidator(engine.State)
	alice := engine.State.Players[0]

	// Add a property card to Alice's hand
	propCard := &PropertyCard{
		BaseCard: BaseCard{ID: "test_prop", Type: CardTypeProperty, Name: "Test", Value: 1},
		Color:    "blue",
	}
	alice.AddToHand(propCard)

	// Valid property play
	action := NewAction(ActionPlayProperty, alice.ID, map[string]string{"cardId": "test_prop"})
	result := validator.Validate(action)
	if !result.Valid {
		t.Errorf("Expected valid, got: %s", result.Error)
	}

	// Card not in hand
	action = NewAction(ActionPlayProperty, alice.ID, map[string]string{"cardId": "nonexistent"})
	result = validator.Validate(action)
	if result.Valid {
		t.Error("Should reject card not in hand")
	}

	// Wrong card type
	moneyCard := &MoneyCard{BaseCard: BaseCard{ID: "test_money", Type: CardTypeMoney, Value: 1}}
	alice.AddToHand(moneyCard)
	action = NewAction(ActionPlayProperty, alice.ID, map[string]string{"cardId": "test_money"})
	result = validator.Validate(action)
	if result.Valid {
		t.Error("Should reject money card as property")
	}
}

func TestValidatorMaxActions(t *testing.T) {
	engine := setupTestGame()
	validator := NewValidator(engine.State)
	alice := engine.State.Players[0]

	// Use up all actions
	engine.State.ActionsPlayedThisTurn = 3

	propCard := &PropertyCard{
		BaseCard: BaseCard{ID: "test_prop", Type: CardTypeProperty, Name: "Test", Value: 1},
		Color:    "blue",
	}
	alice.AddToHand(propCard)

	action := NewAction(ActionPlayProperty, alice.ID, map[string]string{"cardId": "test_prop"})
	result := validator.Validate(action)
	if result.Valid {
		t.Error("Should reject when max actions exceeded")
	}
	if result.Error != "maximum actions per turn exceeded" {
		t.Errorf("Wrong error message: %s", result.Error)
	}
}

func TestValidatorPlayMoney(t *testing.T) {
	engine := setupTestGame()
	validator := NewValidator(engine.State)
	alice := engine.State.Players[0]

	// Add money card
	moneyCard := &MoneyCard{BaseCard: BaseCard{ID: "test_money", Type: CardTypeMoney, Value: 5}}
	alice.AddToHand(moneyCard)

	action := NewAction(ActionPlayMoney, alice.ID, map[string]string{"cardId": "test_money"})
	result := validator.Validate(action)
	if !result.Valid {
		t.Errorf("Expected valid, got: %s", result.Error)
	}

	// Zero value card can't be banked
	zeroCard := &WildcardCard{BaseCard: BaseCard{ID: "zero_wild", Type: CardTypeWildcard, Value: 0}}
	alice.AddToHand(zeroCard)
	action = NewAction(ActionPlayMoney, alice.ID, map[string]string{"cardId": "zero_wild"})
	result = validator.Validate(action)
	if result.Valid {
		t.Error("Should reject $0 card for banking")
	}
}

func TestValidatorPlayRent(t *testing.T) {
	engine := setupTestGame()
	validator := NewValidator(engine.State)
	alice := engine.State.Players[0]

	// Add rent card and property
	rentCard := &RentCard{
		BaseCard: BaseCard{ID: "test_rent", Type: CardTypeRent, Value: 1},
		Colors:   []string{"blue", "brown"},
		Target:   "all",
	}
	alice.AddToHand(rentCard)

	// No properties - should fail for non-wild rent
	action := NewAction(ActionPlayRent, alice.ID, map[string]string{
		"cardId": "test_rent",
		"color":  "blue",
	})
	result := validator.Validate(action)
	if result.Valid {
		t.Error("Should reject rent when no properties owned")
	}

	// Add a blue property
	alice.GetPropertySet("blue").AddCard(&PropertyCard{
		BaseCard: BaseCard{ID: "prop_a"},
		Color:    "blue",
	})

	result = validator.Validate(action)
	if !result.Valid {
		t.Errorf("Expected valid with property, got: %s", result.Error)
	}
}

func TestValidatorPlayRentWild(t *testing.T) {
	engine := setupTestGame()
	validator := NewValidator(engine.State)
	alice := engine.State.Players[0]

	// Wild rent doesn't require properties
	wildRent := &RentCard{
		BaseCard: BaseCard{ID: "wild_rent", Type: CardTypeRent, Name: "Wild Rent", Value: 3},
		Colors:   []string{},
		Target:   "one",
	}
	alice.AddToHand(wildRent)

	action := NewAction(ActionPlayRent, alice.ID, map[string]string{
		"cardId":         "wild_rent",
		"color":          "blue",
		"targetPlayerId": "player_1",
	})

	result := validator.Validate(action)
	if !result.Valid {
		t.Errorf("Wild rent should be valid without properties: %s", result.Error)
	}
}

func TestValidatorEndTurn(t *testing.T) {
	engine := setupTestGame()
	validator := NewValidator(engine.State)
	alice := engine.State.Players[0]
	bob := engine.State.Players[1]

	// Alice can end turn
	action := NewAction(ActionEndTurn, alice.ID, nil)
	result := validator.Validate(action)
	if !result.Valid {
		t.Errorf("Expected valid, got: %s", result.Error)
	}

	// Bob cannot end Alice's turn
	action = NewAction(ActionEndTurn, bob.ID, nil)
	result = validator.Validate(action)
	if result.Valid {
		t.Error("Bob should not be able to end Alice's turn")
	}
}

func TestValidatorFlipWildcard(t *testing.T) {
	engine := setupTestGame()
	validator := NewValidator(engine.State)
	alice := engine.State.Players[0]

	// Add wildcard to properties
	wildcard := &WildcardCard{
		BaseCard:     BaseCard{ID: "test_wild", Type: CardTypeWildcard, Value: 1},
		Colors:       []string{"blue", "brown"},
		CurrentColor: "blue",
	}
	alice.GetPropertySet("blue").AddCard(wildcard)

	// Can flip incomplete set
	action := NewAction(ActionFlipWildcard, alice.ID, map[string]string{"cardId": "test_wild"})
	result := validator.Validate(action)
	if !result.Valid {
		t.Errorf("Expected valid, got: %s", result.Error)
	}

	// Complete the blue set
	alice.GetPropertySet("blue").AddCard(&PropertyCard{BaseCard: BaseCard{ID: "prop_a"}, Color: "blue"})
	alice.GetPropertySet("blue").AddCard(&PropertyCard{BaseCard: BaseCard{ID: "prop_c"}, Color: "blue"})

	// Can't flip in complete set
	result = validator.Validate(action)
	if result.Valid {
		t.Error("Should not allow flipping wildcard in complete set")
	}
}

func TestValidatorUnknownAction(t *testing.T) {
	engine := setupTestGame()
	validator := NewValidator(engine.State)

	action := NewAction("UNKNOWN_ACTION", "player_0", nil)
	result := validator.Validate(action)
	if result.Valid {
		t.Error("Should reject unknown action type")
	}
}
