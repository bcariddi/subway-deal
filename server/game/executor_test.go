package game

import "testing"

func TestExecutorPlayProperty(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]

	propCard := &PropertyCard{
		BaseCard: BaseCard{ID: "test_prop", Type: CardTypeProperty, Name: "A", Value: 1},
		Color:    "blue",
	}
	alice.AddToHand(propCard)

	action := NewAction(ActionPlayProperty, alice.ID, map[string]string{"cardId": "test_prop"})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	blueSet := alice.GetPropertySet("blue")
	if len(blueSet.Cards) != 1 {
		t.Error("Card should be in blue property set")
	}

	if engine.State.ActionsPlayedThisTurn != 1 {
		t.Errorf("Expected 1 action played, got %d", engine.State.ActionsPlayedThisTurn)
	}
}

func TestExecutorPlayWildcard(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]

	wildcard := &WildcardCard{
		BaseCard:     BaseCard{ID: "test_wild", Type: CardTypeWildcard, Name: "Broadway Junction", Value: 1},
		Colors:       []string{"blue", "brown"},
		CurrentColor: "brown",
	}
	alice.AddToHand(wildcard)

	action := NewAction(ActionPlayProperty, alice.ID, map[string]string{"cardId": "test_wild"})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	brownSet := alice.GetPropertySet("brown")
	if len(brownSet.Cards) != 1 {
		t.Error("Wildcard should be in brown property set")
	}
}

func TestExecutorPlayMoney(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]

	moneyCard := &MoneyCard{
		BaseCard: BaseCard{ID: "test_money", Type: CardTypeMoney, Name: "$5", Value: 5},
	}
	alice.AddToHand(moneyCard)

	action := NewAction(ActionPlayMoney, alice.ID, map[string]string{"cardId": "test_money"})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	if alice.GetTotalMoney() != 5 {
		t.Errorf("Expected $5 in bank, got $%d", alice.GetTotalMoney())
	}
}

func TestExecutorPlayRentWithResponse(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]
	bob := engine.State.Players[1]

	// Give Alice a blue property
	alice.GetPropertySet("blue").AddCard(&PropertyCard{
		BaseCard: BaseCard{ID: "prop_a"},
		Color:    "blue",
		Rent:     map[int]int{1: 1, 2: 2, 3: 3},
	})

	// Give Bob some money
	bob.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "bob_money", Value: 5}})

	// Alice plays rent
	rentCard := &RentCard{
		BaseCard: BaseCard{ID: "test_rent", Type: CardTypeRent, Value: 1},
		Colors:   []string{"blue"},
		Target:   "all",
	}
	alice.AddToHand(rentCard)

	action := NewAction(ActionPlayRent, alice.ID, map[string]string{
		"cardId": "test_rent",
		"color":  "blue",
	})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	// Should have pending action
	if !result.PendingAction {
		t.Error("Rent should create pending action")
	}
	if engine.State.TurnPhase != TurnPhaseResponse {
		t.Error("Should be in response phase")
	}

	// Bob accepts (pays)
	acceptAction := NewAction(ActionAccept, bob.ID, nil)
	acceptResult := executor.Execute(acceptAction)

	if !acceptResult.Success {
		t.Errorf("Accept failed: %s", acceptResult.Error)
	}

	// Alice should have received payment
	if alice.GetTotalMoney() < 1 {
		t.Error("Alice should have received rent")
	}

	// Pending action should be cleared
	if engine.State.HasPendingAction() {
		t.Error("Pending action should be cleared after response")
	}
}

func TestExecutorSwipeIn(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]

	initialHandSize := len(alice.Hand)

	swipeIn := &ActionCard{
		BaseCard: BaseCard{ID: "swipe_in", Type: CardTypeAction, Name: "Swipe In", Value: 1},
		Effect:   "Draw 2 extra cards",
	}
	alice.AddToHand(swipeIn)

	action := NewAction(ActionSwipeIn, alice.ID, map[string]string{"cardId": "swipe_in"})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	expectedHand := initialHandSize + 2
	if len(alice.Hand) != expectedHand {
		t.Errorf("Expected %d cards, got %d", expectedHand, len(alice.Hand))
	}
}

func TestExecutorPowerBrokerWithResponse(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]
	bob := engine.State.Players[1]

	// Give Bob a property (incomplete set)
	bobProp := &PropertyCard{
		BaseCard: BaseCard{ID: "bob_prop", Type: CardTypeProperty, Name: "A", Value: 1},
		Color:    "blue",
	}
	bob.GetPropertySet("blue").AddCard(bobProp)

	// Alice plays Power Broker
	powerBroker := &ActionCard{
		BaseCard: BaseCard{ID: "power_broker", Type: CardTypeAction, Name: "Power Broker", Value: 3},
	}
	alice.AddToHand(powerBroker)

	action := NewAction(ActionPowerBroker, alice.ID, map[string]string{
		"cardId":         "power_broker",
		"targetPlayerId": bob.ID,
		"color":          "blue",
		"targetCardId":   "bob_prop",
	})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}
	if !result.PendingAction {
		t.Error("Should create pending action")
	}

	// Bob accepts
	acceptAction := NewAction(ActionAccept, bob.ID, nil)
	executor.Execute(acceptAction)

	// Property should now be in Alice's set
	aliceBlue := alice.GetPropertySet("blue")
	if len(aliceBlue.Cards) != 1 {
		t.Error("Alice should have stolen the property")
	}

	bobBlue := bob.GetPropertySet("blue")
	if len(bobBlue.Cards) != 0 {
		t.Error("Bob should no longer have the property")
	}
}

func TestExecutorPowerBrokerCompleteSet(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]
	bob := engine.State.Players[1]

	// Give Bob a complete set
	bobBlue := bob.GetPropertySet("blue")
	for i := 0; i < 3; i++ {
		bobBlue.AddCard(&PropertyCard{
			BaseCard: BaseCard{ID: "bob_blue_" + string(rune('a'+i))},
			Color:    "blue",
		})
	}

	powerBroker := &ActionCard{
		BaseCard: BaseCard{ID: "power_broker", Type: CardTypeAction, Name: "Power Broker", Value: 3},
	}
	alice.AddToHand(powerBroker)

	action := NewAction(ActionPowerBroker, alice.ID, map[string]string{
		"cardId":         "power_broker",
		"targetPlayerId": bob.ID,
		"color":          "blue",
		"targetCardId":   "bob_blue_a",
	})
	result := executor.Execute(action)

	if result.Success {
		t.Error("Should not be able to steal from complete set")
	}
}

func TestExecutorLineClosureWithResponse(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]
	bob := engine.State.Players[1]

	// Give Bob a complete set with improvements
	bobBlue := bob.GetPropertySet("blue")
	for i := 0; i < 3; i++ {
		bobBlue.AddCard(&PropertyCard{
			BaseCard: BaseCard{ID: "bob_blue_" + string(rune('a'+i))},
			Color:    "blue",
		})
	}
	bobBlue.AddImprovement("express")

	lineClosure := &ActionCard{
		BaseCard: BaseCard{ID: "line_closure", Type: CardTypeAction, Name: "Line Closure", Value: 5},
	}
	alice.AddToHand(lineClosure)

	action := NewAction(ActionLineClosure, alice.ID, map[string]string{
		"cardId":         "line_closure",
		"targetPlayerId": bob.ID,
		"color":          "blue",
	})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	// Bob accepts
	acceptAction := NewAction(ActionAccept, bob.ID, nil)
	executor.Execute(acceptAction)

	aliceBlue := alice.GetPropertySet("blue")
	if len(aliceBlue.Cards) != 3 {
		t.Errorf("Alice should have 3 cards, got %d", len(aliceBlue.Cards))
	}
	if len(aliceBlue.Improvements) != 1 {
		t.Error("Alice should have the improvements")
	}

	if len(bobBlue.Cards) != 0 {
		t.Error("Bob should have no cards")
	}
}

func TestExecutorMissedTrainWithResponse(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]
	bob := engine.State.Players[1]

	bob.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "bob_money", Value: 10}})

	missedTrain := &ActionCard{
		BaseCard: BaseCard{ID: "missed_train", Type: CardTypeAction, Name: "Missed Your Train", Value: 3},
	}
	alice.AddToHand(missedTrain)

	action := NewAction(ActionMissedTrain, alice.ID, map[string]string{
		"cardId":         "missed_train",
		"targetPlayerId": bob.ID,
	})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	// Bob accepts
	acceptAction := NewAction(ActionAccept, bob.ID, nil)
	executor.Execute(acceptAction)

	if alice.GetTotalMoney() < 5 {
		t.Error("Alice should have received $5")
	}
}

func TestExecutorItsMyStopWithResponse(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob", "Charlie"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]
	bob := engine.State.Players[1]
	charlie := engine.State.Players[2]

	bob.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "bob_money_1", Value: 1}})
	bob.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "bob_money_2", Value: 1}})
	charlie.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "charlie_money_1", Value: 1}})
	charlie.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "charlie_money_2", Value: 1}})

	itsMyStop := &ActionCard{
		BaseCard: BaseCard{ID: "its_my_stop", Type: CardTypeAction, Name: "It's My Stop!", Value: 2},
	}
	alice.AddToHand(itsMyStop)

	action := NewAction(ActionItsMyStop, alice.ID, map[string]string{"cardId": "its_my_stop"})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	// Both players accept
	executor.Execute(NewAction(ActionAccept, bob.ID, nil))
	executor.Execute(NewAction(ActionAccept, charlie.ID, nil))

	if alice.GetTotalMoney() != 4 {
		t.Errorf("Expected $4 from two players, got $%d", alice.GetTotalMoney())
	}
}

func TestExecutorExpressService(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]

	aliceBlue := alice.GetPropertySet("blue")
	for i := 0; i < 3; i++ {
		aliceBlue.AddCard(&PropertyCard{
			BaseCard: BaseCard{ID: "alice_blue_" + string(rune('a'+i))},
			Color:    "blue",
			Rent:     map[int]int{1: 1, 2: 2, 3: 3},
		})
	}

	initialRent := aliceBlue.GetRent()

	expressService := &ActionCard{
		BaseCard: BaseCard{ID: "express_service", Type: CardTypeAction, Name: "Express Service", Value: 3},
	}
	alice.AddToHand(expressService)

	action := NewAction(ActionExpressService, alice.ID, map[string]string{
		"cardId": "express_service",
		"color":  "blue",
	})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	if aliceBlue.GetRent() != initialRent+3 {
		t.Errorf("Expected rent to increase by $3, got %d -> %d", initialRent, aliceBlue.GetRent())
	}
}

func TestExecutorFlipWildcard(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]

	wildcard := &WildcardCard{
		BaseCard:     BaseCard{ID: "test_wild", Type: CardTypeWildcard, Name: "Broadway Junction", Value: 1},
		Colors:       []string{"blue", "brown"},
		CurrentColor: "blue",
	}
	alice.GetPropertySet("blue").AddCard(wildcard)

	action := NewAction(ActionFlipWildcard, alice.ID, map[string]string{"cardId": "test_wild"})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	brownSet := alice.GetPropertySet("brown")
	if len(brownSet.Cards) != 1 {
		t.Error("Wildcard should be in brown set")
	}

	blueSet := alice.GetPropertySet("blue")
	if len(blueSet.Cards) != 0 {
		t.Error("Wildcard should no longer be in blue set")
	}
}

func TestExecutorFareEvasion(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]
	bob := engine.State.Players[1]

	// Give Alice a property
	alice.GetPropertySet("blue").AddCard(&PropertyCard{
		BaseCard: BaseCard{ID: "prop_a"},
		Color:    "blue",
		Rent:     map[int]int{1: 3},
	})

	// Give Bob money and Fare Evasion
	bob.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "bob_money", Value: 10}})
	fareEvasion := &ActionCard{
		BaseCard: BaseCard{ID: "fare_evasion", Type: CardTypeAction, Name: "Fare Evasion", Value: 4},
	}
	bob.AddToHand(fareEvasion)

	// Alice plays rent
	rentCard := &RentCard{
		BaseCard: BaseCard{ID: "test_rent", Type: CardTypeRent, Value: 1},
		Colors:   []string{"blue"},
		Target:   "all",
	}
	alice.AddToHand(rentCard)

	action := NewAction(ActionPlayRent, alice.ID, map[string]string{
		"cardId": "test_rent",
		"color":  "blue",
	})
	executor.Execute(action)

	// Bob plays Fare Evasion
	fareEvasionAction := NewAction(ActionPlayFareEvasion, bob.ID, map[string]string{
		"cardId": "fare_evasion",
	})
	result := executor.Execute(fareEvasionAction)

	if !result.Success {
		t.Errorf("Fare Evasion failed: %s", result.Error)
	}

	// Alice should NOT have received payment
	if alice.GetTotalMoney() != 0 {
		t.Errorf("Alice should not have received rent, got $%d", alice.GetTotalMoney())
	}

	// Bob should still have money
	if bob.GetTotalMoney() != 10 {
		t.Errorf("Bob should still have $10, got $%d", bob.GetTotalMoney())
	}

	// Pending action should be cleared
	if engine.State.HasPendingAction() {
		t.Error("Pending action should be cleared after Fare Evasion")
	}
}

func TestExecutorRushHour(t *testing.T) {
	engine := NewEngine("test-game", []string{"Alice", "Bob"})
	engine.StartTurn()
	executor := NewExecutor(engine.State)
	alice := engine.State.Players[0]
	bob := engine.State.Players[1]

	// Give Alice a blue property with $3 rent
	alice.GetPropertySet("blue").AddCard(&PropertyCard{
		BaseCard: BaseCard{ID: "prop_a"},
		Color:    "blue",
		Rent:     map[int]int{1: 3},
	})

	// Give Bob exact change for $6
	for i := 0; i < 6; i++ {
		bob.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "bob_money_" + string(rune('a'+i)), Value: 1}})
	}

	// Alice has rent card and Rush Hour
	rentCard := &RentCard{
		BaseCard: BaseCard{ID: "test_rent", Type: CardTypeRent, Value: 1},
		Colors:   []string{"blue"},
		Target:   "all",
	}
	rushHour := &ActionCard{
		BaseCard: BaseCard{ID: "rush_hour", Type: CardTypeAction, Name: "Rush Hour", Value: 1},
	}
	alice.AddToHand(rentCard)
	alice.AddToHand(rushHour)

	// Play rent with Rush Hour
	action := NewAction(ActionPlayRent, alice.ID, map[string]string{
		"cardId":        "test_rent",
		"color":         "blue",
		"rushHourCardId": "rush_hour",
	})
	result := executor.Execute(action)

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	// Rent should be doubled
	if engine.State.PendingAction.RentAmount != 6 {
		t.Errorf("Rush Hour should double rent to $6, got $%d", engine.State.PendingAction.RentAmount)
	}

	// Bob accepts
	executor.Execute(NewAction(ActionAccept, bob.ID, nil))

	// Alice should receive $6
	if alice.GetTotalMoney() != 6 {
		t.Errorf("Alice should have received $6, got $%d", alice.GetTotalMoney())
	}
}
