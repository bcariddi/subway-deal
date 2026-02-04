package game

import (
	"errors"
	"sort"
)

// Player represents a game participant
type Player struct {
	ID         string                  `json:"id"`
	Name       string                  `json:"name"`
	Hand       []Card                  `json:"hand"`
	Bank       []Card                  `json:"bank"`
	Properties map[string]*PropertySet `json:"properties"`
}

func NewPlayer(id, name string) *Player {
	return &Player{
		ID:         id,
		Name:       name,
		Hand:       make([]Card, 0),
		Bank:       make([]Card, 0),
		Properties: make(map[string]*PropertySet),
	}
}

func (p *Player) AddToHand(card Card) {
	p.Hand = append(p.Hand, card)
}

func (p *Player) RemoveFromHand(cardID string) (Card, error) {
	for i, card := range p.Hand {
		if card.GetID() == cardID {
			removed := p.Hand[i]
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			return removed, nil
		}
	}
	return nil, errors.New("card not in hand")
}

func (p *Player) GetCardFromHand(cardID string) Card {
	for _, card := range p.Hand {
		if card.GetID() == cardID {
			return card
		}
	}
	return nil
}

func (p *Player) AddToBank(card Card) {
	p.Bank = append(p.Bank, card)
}

func (p *Player) GetTotalMoney() int {
	total := 0
	for _, card := range p.Bank {
		total += card.GetValue()
	}
	return total
}

// PayMoney removes cards from bank to pay an amount
// Returns the cards paid and total value (may exceed amount)
func (p *Player) PayMoney(amount int) ([]Card, int) {
	if amount <= 0 {
		return nil, 0
	}

	// Sort bank by value (ascending) to minimize overpayment
	sorted := make([]Card, len(p.Bank))
	copy(sorted, p.Bank)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].GetValue() < sorted[j].GetValue()
	})

	payment := make([]Card, 0)
	total := 0

	for _, card := range sorted {
		if total >= amount {
			break
		}
		payment = append(payment, card)
		total += card.GetValue()

		// Remove from bank
		for i, bankCard := range p.Bank {
			if bankCard.GetID() == card.GetID() {
				p.Bank = append(p.Bank[:i], p.Bank[i+1:]...)
				break
			}
		}
	}

	return payment, total
}

func (p *Player) GetPropertySet(color string) *PropertySet {
	if ps, ok := p.Properties[color]; ok {
		return ps
	}
	ps := NewPropertySet(color)
	p.Properties[color] = ps
	return ps
}

func (p *Player) GetCompleteSetCount() int {
	count := 0
	for _, ps := range p.Properties {
		if ps.IsComplete() {
			count++
		}
	}
	return count
}

func (p *Player) HasWon() bool {
	return p.GetCompleteSetCount() >= 3
}
