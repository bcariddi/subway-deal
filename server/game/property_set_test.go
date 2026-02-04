package game

import (
	"fmt"
	"testing"
)

func TestPropertySetCompletion(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		cards    int
		expected bool
	}{
		{"brown incomplete", "brown", 1, false},
		{"brown complete", "brown", 2, true},
		{"blue incomplete", "blue", 2, false},
		{"blue complete", "blue", 3, true},
		{"railroad incomplete", "railroad", 3, false},
		{"railroad complete", "railroad", 4, true},
		{"utility incomplete", "utility", 1, false},
		{"utility complete", "utility", 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := NewPropertySet(tt.color)
			for i := 0; i < tt.cards; i++ {
				ps.AddCard(&PropertyCard{
					BaseCard: BaseCard{ID: fmt.Sprintf("test_%d", i)},
					Color:    tt.color,
				})
			}

			if ps.IsComplete() != tt.expected {
				t.Errorf("IsComplete() = %v, want %v", ps.IsComplete(), tt.expected)
			}
		})
	}
}

func TestPropertySetSize(t *testing.T) {
	tests := []struct {
		color    string
		expected int
	}{
		{"brown", 2},
		{"blue", 3},
		{"pink", 3},
		{"orange", 3},
		{"red", 3},
		{"yellow", 3},
		{"green", 3},
		{"darkblue", 2},
		{"railroad", 4},
		{"utility", 2},
		{"unknown", 3}, // defaults to 3
	}

	for _, tt := range tests {
		t.Run(tt.color, func(t *testing.T) {
			ps := NewPropertySet(tt.color)
			if ps.SetSize() != tt.expected {
				t.Errorf("SetSize() = %v, want %v", ps.SetSize(), tt.expected)
			}
		})
	}
}

func TestPropertySetAddRemoveCard(t *testing.T) {
	ps := NewPropertySet("blue")

	card1 := &PropertyCard{BaseCard: BaseCard{ID: "card_1", Name: "A"}, Color: "blue"}
	card2 := &PropertyCard{BaseCard: BaseCard{ID: "card_2", Name: "C"}, Color: "blue"}

	ps.AddCard(card1)
	ps.AddCard(card2)

	if len(ps.Cards) != 2 {
		t.Errorf("Expected 2 cards, got %d", len(ps.Cards))
	}

	removed, err := ps.RemoveCard("card_1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if removed.GetID() != "card_1" {
		t.Errorf("Removed wrong card: got %s", removed.GetID())
	}
	if len(ps.Cards) != 1 {
		t.Errorf("Expected 1 card after removal, got %d", len(ps.Cards))
	}

	// Try to remove non-existent card
	_, err = ps.RemoveCard("nonexistent")
	if err == nil {
		t.Error("Expected error when removing non-existent card")
	}
}

func TestPropertySetImprovements(t *testing.T) {
	ps := NewPropertySet("blue")

	// Can't add improvement to incomplete set
	if ps.CanAddImprovement("express") {
		t.Error("Should not be able to add improvement to incomplete set")
	}

	// Complete the set
	for i := 0; i < 3; i++ {
		ps.AddCard(&PropertyCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("card_%d", i)},
			Color:    "blue",
		})
	}

	// Now can add express
	if !ps.CanAddImprovement("express") {
		t.Error("Should be able to add express to complete set")
	}

	// Can't add station before express
	if ps.CanAddImprovement("station") {
		t.Error("Should not be able to add station before express")
	}

	// Add express
	err := ps.AddImprovement("express")
	if err != nil {
		t.Errorf("Unexpected error adding express: %v", err)
	}

	// Can't add express twice
	if ps.CanAddImprovement("express") {
		t.Error("Should not be able to add express twice")
	}

	// Now can add station
	if !ps.CanAddImprovement("station") {
		t.Error("Should be able to add station after express")
	}

	err = ps.AddImprovement("station")
	if err != nil {
		t.Errorf("Unexpected error adding station: %v", err)
	}

	// Can't add station twice
	if ps.CanAddImprovement("station") {
		t.Error("Should not be able to add station twice")
	}
}

func TestPropertySetImprovementsRailroadUtility(t *testing.T) {
	// Railroads can't have improvements
	railroad := NewPropertySet("railroad")
	for i := 0; i < 4; i++ {
		railroad.AddCard(&PropertyCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("rr_%d", i)},
			Color:    "railroad",
		})
	}
	if railroad.CanAddImprovement("express") {
		t.Error("Railroad should not allow improvements")
	}

	// Utilities can't have improvements
	utility := NewPropertySet("utility")
	for i := 0; i < 2; i++ {
		utility.AddCard(&PropertyCard{
			BaseCard: BaseCard{ID: fmt.Sprintf("util_%d", i)},
			Color:    "utility",
		})
	}
	if utility.CanAddImprovement("express") {
		t.Error("Utility should not allow improvements")
	}
}

func TestPropertySetGetRent(t *testing.T) {
	ps := NewPropertySet("blue")

	// Empty set has no rent
	if ps.GetRent() != 0 {
		t.Errorf("Empty set should have 0 rent, got %d", ps.GetRent())
	}

	// Add first property
	ps.AddCard(&PropertyCard{
		BaseCard: BaseCard{ID: "card_1"},
		Color:    "blue",
		Rent:     map[int]int{1: 1, 2: 2, 3: 3},
	})
	if ps.GetRent() != 1 {
		t.Errorf("1 property should have $1 rent, got %d", ps.GetRent())
	}

	// Add second property
	ps.AddCard(&PropertyCard{
		BaseCard: BaseCard{ID: "card_2"},
		Color:    "blue",
		Rent:     map[int]int{1: 1, 2: 2, 3: 3},
	})
	if ps.GetRent() != 2 {
		t.Errorf("2 properties should have $2 rent, got %d", ps.GetRent())
	}

	// Complete set
	ps.AddCard(&PropertyCard{
		BaseCard: BaseCard{ID: "card_3"},
		Color:    "blue",
		Rent:     map[int]int{1: 1, 2: 2, 3: 3},
	})
	if ps.GetRent() != 3 {
		t.Errorf("3 properties should have $3 rent, got %d", ps.GetRent())
	}

	// Add express service (+$3)
	ps.AddImprovement("express")
	if ps.GetRent() != 6 {
		t.Errorf("With express should have $6 rent, got %d", ps.GetRent())
	}

	// Add station (+$4)
	ps.AddImprovement("station")
	if ps.GetRent() != 10 {
		t.Errorf("With station should have $10 rent, got %d", ps.GetRent())
	}
}
