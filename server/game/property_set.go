package game

import "errors"

// SetSizes defines how many cards complete each color set
var SetSizes = map[string]int{
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

// PropertySet tracks cards of one color owned by a player
type PropertySet struct {
	Color        string   `json:"color"`
	Cards        []Card   `json:"cards"`
	Improvements []string `json:"improvements"` // "express", "station"
}

func NewPropertySet(color string) *PropertySet {
	return &PropertySet{
		Color:        color,
		Cards:        make([]Card, 0),
		Improvements: make([]string, 0),
	}
}

func (ps *PropertySet) SetSize() int {
	if size, ok := SetSizes[ps.Color]; ok {
		return size
	}
	return 3
}

func (ps *PropertySet) AddCard(card Card) {
	ps.Cards = append(ps.Cards, card)
}

func (ps *PropertySet) RemoveCard(cardID string) (Card, error) {
	for i, card := range ps.Cards {
		if card.GetID() == cardID {
			removed := ps.Cards[i]
			ps.Cards = append(ps.Cards[:i], ps.Cards[i+1:]...)
			return removed, nil
		}
	}
	return nil, errors.New("card not found in set")
}

func (ps *PropertySet) IsComplete() bool {
	return len(ps.Cards) >= ps.SetSize()
}

func (ps *PropertySet) CanAddImprovement(improvementType string) bool {
	if !ps.IsComplete() {
		return false
	}

	// Railroads and utilities can't have improvements
	if ps.Color == "railroad" || ps.Color == "utility" {
		return false
	}

	switch improvementType {
	case "express":
		return !ps.hasImprovement("express")
	case "station":
		return ps.hasImprovement("express") && !ps.hasImprovement("station")
	}
	return false
}

func (ps *PropertySet) hasImprovement(t string) bool {
	for _, imp := range ps.Improvements {
		if imp == t {
			return true
		}
	}
	return false
}

func (ps *PropertySet) AddImprovement(improvementType string) error {
	if !ps.CanAddImprovement(improvementType) {
		return errors.New("cannot add improvement")
	}
	ps.Improvements = append(ps.Improvements, improvementType)
	return nil
}

func (ps *PropertySet) GetRent() int {
	if len(ps.Cards) == 0 {
		return 0
	}

	// Get base rent from first property card
	var baseRent int
	for _, card := range ps.Cards {
		if prop, ok := card.(*PropertyCard); ok {
			baseRent = prop.GetRent(len(ps.Cards))
			break
		}
	}

	// Add improvement bonuses
	bonus := 0
	if ps.hasImprovement("express") {
		bonus += 3
	}
	if ps.hasImprovement("station") {
		bonus += 4
	}

	return baseRent + bonus
}
