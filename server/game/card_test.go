package game

import (
	"testing"
)

func TestBaseCard(t *testing.T) {
	bc := BaseCard{
		ID:    "test-1",
		Type:  CardTypeMoney,
		Name:  "Test Card",
		Value: 1,
	}

	if bc.GetID() != "test-1" {
		t.Errorf("Expected ID test-1, got %s", bc.GetID())
	}
	if bc.GetType() != CardTypeMoney {
		t.Errorf("Expected type money, got %s", bc.GetType())
	}
	if bc.GetName() != "Test Card" {
		t.Errorf("Expected name Test Card, got %s", bc.GetName())
	}
	if bc.GetValue() != 1 {
		t.Errorf("Expected value 1, got %d", bc.GetValue())
	}
}

func TestPropertyCard(t *testing.T) {
	pc := PropertyCard{
		BaseCard: BaseCard{ID: "prop-a", Type: CardTypeProperty, Name: "A", Value: 1},
		Color:    "blue",
		ColorHex: "#0039A6",
		SetSize:  3,
		Position: 1,
		Rent:     map[int]int{1: 1, 2: 2, 3: 3},
	}

	// Test interface methods
	var card Card = &pc
	if card.GetID() != "prop-a" {
		t.Error("PropertyCard should implement Card interface")
	}

	// Test rent calculation
	if pc.GetRent(1) != 1 {
		t.Errorf("Expected $1 rent for 1 property, got %d", pc.GetRent(1))
	}
	if pc.GetRent(2) != 2 {
		t.Errorf("Expected $2 rent for 2 properties, got %d", pc.GetRent(2))
	}
	if pc.GetRent(3) != 3 {
		t.Errorf("Expected $3 rent for 3 properties, got %d", pc.GetRent(3))
	}
	if pc.GetRent(5) != 0 {
		t.Errorf("Expected $0 rent for invalid count, got %d", pc.GetRent(5))
	}
}

func TestWildcardCard(t *testing.T) {
	wc := WildcardCard{
		BaseCard:     BaseCard{ID: "wild-1", Type: CardTypeWildcard, Name: "Broadway Junction", Value: 1},
		Colors:       []string{"blue", "brown"},
		CurrentColor: "blue",
		Description:  "Can be used as Light Blue or Brown",
	}

	// Test interface
	var card Card = &wc
	if card.GetType() != CardTypeWildcard {
		t.Error("WildcardCard should implement Card interface")
	}

	// Test CanBeColor
	if !wc.CanBeColor("blue") {
		t.Error("Should be able to be blue")
	}
	if !wc.CanBeColor("brown") {
		t.Error("Should be able to be brown")
	}
	if wc.CanBeColor("red") {
		t.Error("Should not be able to be red")
	}

	// Test FlipColor
	if wc.CurrentColor != "blue" {
		t.Errorf("Initial color should be blue, got %s", wc.CurrentColor)
	}

	wc.FlipColor()
	if wc.CurrentColor != "brown" {
		t.Errorf("After flip should be brown, got %s", wc.CurrentColor)
	}

	wc.FlipColor()
	if wc.CurrentColor != "blue" {
		t.Errorf("After second flip should be blue, got %s", wc.CurrentColor)
	}
}

func TestWildcardFlipWithInvalidCurrentColor(t *testing.T) {
	wc := WildcardCard{
		BaseCard:     BaseCard{ID: "wild-1"},
		Colors:       []string{"blue", "brown"},
		CurrentColor: "invalid", // Not in Colors
	}

	// FlipColor should not panic with invalid current color
	wc.FlipColor()
	// Color remains unchanged since it wasn't in the list
	if wc.CurrentColor != "invalid" {
		t.Errorf("Color should remain unchanged, got %s", wc.CurrentColor)
	}
}

func TestActionCard(t *testing.T) {
	ac := ActionCard{
		BaseCard: BaseCard{ID: "action-1", Type: CardTypeAction, Name: "Swipe In", Value: 1},
		Effect:   "Draw 2 extra cards",
		MTATheme: "Pass Go",
	}

	var card Card = &ac
	if card.GetType() != CardTypeAction {
		t.Error("ActionCard should implement Card interface")
	}
	if ac.Effect != "Draw 2 extra cards" {
		t.Error("Effect not set correctly")
	}
}

func TestRentCard(t *testing.T) {
	// Standard rent card
	rc := RentCard{
		BaseCard:   BaseCard{ID: "rent-1", Type: CardTypeRent, Name: "Rent", Value: 1},
		Colors:     []string{"blue", "brown"},
		ColorNames: []string{"Light Blue", "Brown"},
		Target:     "all",
	}

	var card Card = &rc
	if card.GetType() != CardTypeRent {
		t.Error("RentCard should implement Card interface")
	}

	if rc.IsWildRent() {
		t.Error("Standard rent should not be wild")
	}

	// Wild rent card
	wildRent := RentCard{
		BaseCard: BaseCard{ID: "rent-wild", Name: "Wild Rent", Value: 3},
		Colors:   []string{},
		Target:   "one",
	}

	if !wildRent.IsWildRent() {
		t.Error("Wild rent should be identified as wild")
	}

	// Also test by name
	namedWild := RentCard{
		BaseCard: BaseCard{Name: "Wild Rent"},
		Colors:   []string{"blue"}, // Has colors but name is Wild Rent
	}
	if !namedWild.IsWildRent() {
		t.Error("Card named 'Wild Rent' should be wild")
	}
}

func TestMoneyCard(t *testing.T) {
	mc := MoneyCard{
		BaseCard:     BaseCard{ID: "money-1", Type: CardTypeMoney, Name: "$1", Value: 1},
		Denomination: 1,
		DisplayValue: "$1",
		Theme:        "Base Fare",
	}

	var card Card = &mc
	if card.GetType() != CardTypeMoney {
		t.Error("MoneyCard should implement Card interface")
	}
	if card.GetValue() != 1 {
		t.Errorf("Expected value 1, got %d", card.GetValue())
	}
}

func TestCardTypes(t *testing.T) {
	// Verify all card type constants
	types := []CardType{
		CardTypeProperty,
		CardTypeWildcard,
		CardTypeAction,
		CardTypeRent,
		CardTypeMoney,
	}

	expected := []string{"property", "wildcard", "action", "rent", "money"}

	for i, ct := range types {
		if string(ct) != expected[i] {
			t.Errorf("CardType %d: expected %s, got %s", i, expected[i], ct)
		}
	}
}
