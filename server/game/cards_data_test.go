package game

import "testing"

func TestCreateAllCards(t *testing.T) {
	cards := CreateAllCards()

	if len(cards) != 106 {
		t.Errorf("Expected 106 cards, got %d", len(cards))
	}

	// Count by type
	counts := map[CardType]int{}
	for _, card := range cards {
		counts[card.GetType()]++
	}

	// Verify counts match spec
	expected := map[CardType]int{
		CardTypeProperty: 28,
		CardTypeWildcard: 11,
		CardTypeAction:   34,
		CardTypeRent:     13,
		CardTypeMoney:    20,
	}

	for cardType, expectedCount := range expected {
		if counts[cardType] != expectedCount {
			t.Errorf("Expected %d %s cards, got %d", expectedCount, cardType, counts[cardType])
		}
	}
}

func TestCreatePropertyCards(t *testing.T) {
	cards := createPropertyCards()

	if len(cards) != 28 {
		t.Errorf("Expected 28 property cards, got %d", len(cards))
	}

	// Count by color
	colorCounts := map[string]int{}
	for _, card := range cards {
		if prop, ok := card.(*PropertyCard); ok {
			colorCounts[prop.Color]++
		}
	}

	expectedColors := map[string]int{
		"brown":    2,
		"blue":     3,
		"pink":     3,
		"orange":   3,
		"red":      3,
		"yellow":   3,
		"green":    3,
		"darkblue": 2,
		"railroad": 4,
		"utility":  2,
	}

	for color, expected := range expectedColors {
		if colorCounts[color] != expected {
			t.Errorf("Expected %d %s properties, got %d", expected, color, colorCounts[color])
		}
	}
}

func TestCreateWildcardCards(t *testing.T) {
	cards := createWildcardCards()

	if len(cards) != 11 {
		t.Errorf("Expected 11 wildcard cards, got %d", len(cards))
	}

	// Verify each wildcard has at least one color
	for _, card := range cards {
		wc := card.(*WildcardCard)
		if len(wc.Colors) == 0 {
			t.Errorf("Wildcard %s has no colors", wc.Name)
		}
		if wc.CurrentColor == "" {
			t.Errorf("Wildcard %s has no current color", wc.Name)
		}
	}

	// Check Fulton Center has all colors and $0 value
	fultonCount := 0
	for _, card := range cards {
		wc := card.(*WildcardCard)
		if wc.Name == "Fulton Center" {
			fultonCount++
			if wc.Value != 0 {
				t.Errorf("Fulton Center should have $0 value, got %d", wc.Value)
			}
			if len(wc.Colors) != 10 {
				t.Errorf("Fulton Center should have 10 colors, got %d", len(wc.Colors))
			}
		}
	}
	if fultonCount != 2 {
		t.Errorf("Expected 2 Fulton Center cards, got %d", fultonCount)
	}
}

func TestCreateActionCards(t *testing.T) {
	cards := createActionCards()

	if len(cards) != 34 {
		t.Errorf("Expected 34 action cards, got %d", len(cards))
	}

	// Count by name
	nameCounts := map[string]int{}
	for _, card := range cards {
		nameCounts[card.GetName()]++
	}

	expectedCounts := map[string]int{
		"Swipe In":        10,
		"Fare Evasion":    3,
		"Power Broker":    3,
		"Service Change":  3,
		"Line Closure":    2,
		"Missed Your Train": 3,
		"It's My Stop!":   3,
		"Rush Hour":       2,
		"Express Service": 3,
		"New Station":     2,
	}

	for name, expected := range expectedCounts {
		if nameCounts[name] != expected {
			t.Errorf("Expected %d '%s' cards, got %d", expected, name, nameCounts[name])
		}
	}
}

func TestCreateRentCards(t *testing.T) {
	cards := createRentCards()

	if len(cards) != 13 {
		t.Errorf("Expected 13 rent cards, got %d", len(cards))
	}

	// Count wild rent cards
	wildCount := 0
	for _, card := range cards {
		rc := card.(*RentCard)
		if rc.IsWildRent() {
			wildCount++
			if rc.Target != "one" {
				t.Error("Wild rent should target one player")
			}
		} else {
			if rc.Target != "all" {
				t.Error("Standard rent should target all players")
			}
		}
	}

	if wildCount != 3 {
		t.Errorf("Expected 3 wild rent cards, got %d", wildCount)
	}
}

func TestCreateMoneyCards(t *testing.T) {
	cards := createMoneyCards()

	if len(cards) != 20 {
		t.Errorf("Expected 20 money cards, got %d", len(cards))
	}

	// Count by denomination and calculate total
	denomCounts := map[int]int{}
	totalValue := 0
	for _, card := range cards {
		mc := card.(*MoneyCard)
		denomCounts[mc.Denomination]++
		totalValue += mc.Value
	}

	expectedDenoms := map[int]int{
		1:  6,
		2:  5,
		3:  3,
		4:  3,
		5:  2,
		10: 1,
	}

	for denom, expected := range expectedDenoms {
		if denomCounts[denom] != expected {
			t.Errorf("Expected %d $%d cards, got %d", expected, denom, denomCounts[denom])
		}
	}

	// Total value: 6*1 + 5*2 + 3*3 + 3*4 + 2*5 + 1*10 = 6+10+9+12+10+10 = 57
	if totalValue != 57 {
		t.Errorf("Expected total money value of $57, got $%d", totalValue)
	}
}

func TestCardIDsAreUnique(t *testing.T) {
	cards := CreateAllCards()
	seen := make(map[string]bool)

	for _, card := range cards {
		id := card.GetID()
		if seen[id] {
			t.Errorf("Duplicate card ID: %s", id)
		}
		seen[id] = true
	}
}
