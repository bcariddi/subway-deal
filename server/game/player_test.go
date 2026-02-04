package game

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	p := NewPlayer("p1", "Alice")

	if p.ID != "p1" {
		t.Errorf("Expected ID p1, got %s", p.ID)
	}
	if p.Name != "Alice" {
		t.Errorf("Expected Name Alice, got %s", p.Name)
	}
	if len(p.Hand) != 0 {
		t.Errorf("Expected empty hand, got %d cards", len(p.Hand))
	}
	if len(p.Bank) != 0 {
		t.Errorf("Expected empty bank, got %d cards", len(p.Bank))
	}
	if len(p.Properties) != 0 {
		t.Errorf("Expected empty properties, got %d sets", len(p.Properties))
	}
}

func TestPlayerHandOperations(t *testing.T) {
	p := NewPlayer("p1", "Alice")

	card1 := &MoneyCard{BaseCard: BaseCard{ID: "m1", Name: "$1", Value: 1}}
	card2 := &MoneyCard{BaseCard: BaseCard{ID: "m2", Name: "$2", Value: 2}}

	// Add to hand
	p.AddToHand(card1)
	p.AddToHand(card2)

	if len(p.Hand) != 2 {
		t.Errorf("Expected 2 cards in hand, got %d", len(p.Hand))
	}

	// Get card from hand
	found := p.GetCardFromHand("m1")
	if found == nil || found.GetID() != "m1" {
		t.Error("Failed to get card from hand")
	}

	// Get non-existent card
	notFound := p.GetCardFromHand("nonexistent")
	if notFound != nil {
		t.Error("Should return nil for non-existent card")
	}

	// Remove from hand
	removed, err := p.RemoveFromHand("m1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if removed.GetID() != "m1" {
		t.Errorf("Removed wrong card: %s", removed.GetID())
	}
	if len(p.Hand) != 1 {
		t.Errorf("Expected 1 card after removal, got %d", len(p.Hand))
	}

	// Remove non-existent card
	_, err = p.RemoveFromHand("nonexistent")
	if err == nil {
		t.Error("Expected error when removing non-existent card")
	}
}

func TestPlayerBankOperations(t *testing.T) {
	p := NewPlayer("p1", "Alice")

	card1 := &MoneyCard{BaseCard: BaseCard{ID: "m1", Name: "$1", Value: 1}}
	card2 := &MoneyCard{BaseCard: BaseCard{ID: "m2", Name: "$5", Value: 5}}

	p.AddToBank(card1)
	p.AddToBank(card2)

	if len(p.Bank) != 2 {
		t.Errorf("Expected 2 cards in bank, got %d", len(p.Bank))
	}

	total := p.GetTotalMoney()
	if total != 6 {
		t.Errorf("Expected $6 total, got %d", total)
	}
}

func TestPlayerPayMoney(t *testing.T) {
	p := NewPlayer("p1", "Alice")

	// Add money to bank
	p.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "m1", Value: 1}})
	p.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "m2", Value: 2}})
	p.AddToBank(&MoneyCard{BaseCard: BaseCard{ID: "m3", Value: 5}})

	// Pay $3 - should use $1 + $2 cards
	payment, total := p.PayMoney(3)
	if total != 3 {
		t.Errorf("Expected to pay $3, paid %d", total)
	}
	if len(payment) != 2 {
		t.Errorf("Expected 2 cards for payment, got %d", len(payment))
	}
	if len(p.Bank) != 1 {
		t.Errorf("Expected 1 card left in bank, got %d", len(p.Bank))
	}

	// Pay $0 - should return nil
	payment, total = p.PayMoney(0)
	if payment != nil || total != 0 {
		t.Error("Paying $0 should return nil, 0")
	}

	// Pay more than available - pays what's available
	payment, total = p.PayMoney(10)
	if total != 5 {
		t.Errorf("Expected to pay $5 (all remaining), paid %d", total)
	}
	if len(p.Bank) != 0 {
		t.Errorf("Expected empty bank after paying everything, got %d", len(p.Bank))
	}
}

func TestPlayerPropertySets(t *testing.T) {
	p := NewPlayer("p1", "Alice")

	// Get property set (creates if doesn't exist)
	blueSet := p.GetPropertySet("blue")
	if blueSet == nil {
		t.Fatal("GetPropertySet should return a set")
	}
	if blueSet.Color != "blue" {
		t.Errorf("Expected blue color, got %s", blueSet.Color)
	}

	// Get same set again
	blueSet2 := p.GetPropertySet("blue")
	if blueSet != blueSet2 {
		t.Error("GetPropertySet should return same instance")
	}

	// Complete sets count
	if p.GetCompleteSetCount() != 0 {
		t.Error("Should have 0 complete sets initially")
	}

	// Add cards to complete blue set
	for i := 0; i < 3; i++ {
		blueSet.AddCard(&PropertyCard{
			BaseCard: BaseCard{ID: "b" + string(rune('1'+i))},
			Color:    "blue",
		})
	}

	if p.GetCompleteSetCount() != 1 {
		t.Errorf("Should have 1 complete set, got %d", p.GetCompleteSetCount())
	}

	// Win condition
	if p.HasWon() {
		t.Error("Should not have won with 1 complete set")
	}

	// Complete 2 more sets
	brownSet := p.GetPropertySet("brown")
	for i := 0; i < 2; i++ {
		brownSet.AddCard(&PropertyCard{
			BaseCard: BaseCard{ID: "br" + string(rune('1'+i))},
			Color:    "brown",
		})
	}

	utilitySet := p.GetPropertySet("utility")
	for i := 0; i < 2; i++ {
		utilitySet.AddCard(&PropertyCard{
			BaseCard: BaseCard{ID: "u" + string(rune('1'+i))},
			Color:    "utility",
		})
	}

	if p.GetCompleteSetCount() != 3 {
		t.Errorf("Should have 3 complete sets, got %d", p.GetCompleteSetCount())
	}

	if !p.HasWon() {
		t.Error("Should have won with 3 complete sets")
	}
}
