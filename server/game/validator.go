package game

import "fmt"

// ValidationResult contains the result of action validation
type ValidationResult struct {
	Valid bool
	Error string
}

func Valid() ValidationResult {
	return ValidationResult{Valid: true}
}

func Invalid(err string) ValidationResult {
	return ValidationResult{Valid: false, Error: err}
}

// Validator validates game actions
type Validator struct {
	state *GameState
}

func NewValidator(state *GameState) *Validator {
	return &Validator{state: state}
}

// Validate checks if an action is valid
func (v *Validator) Validate(action *Action) ValidationResult {
	// Response actions can only happen during response phase
	if action.Type == ActionAccept || action.Type == ActionPlayFareEvasion {
		return v.validateResponse(action)
	}

	// Can't do normal actions during response phase
	if v.state.TurnPhase == TurnPhaseResponse {
		return Invalid("waiting for response to pending action")
	}

	switch action.Type {
	case ActionDrawCards:
		return v.validateDrawCards(action)
	case ActionPlayProperty:
		return v.validatePlayProperty(action)
	case ActionPlayMoney:
		return v.validatePlayMoney(action)
	case ActionPlayRent:
		return v.validatePlayRent(action)
	case ActionPlayAction, ActionSwipeIn, ActionPowerBroker, ActionLineClosure,
		ActionServiceChange, ActionMissedTrain, ActionItsMyStop,
		ActionRushHour, ActionExpressService, ActionNewStation:
		return v.validatePlayAction(action)
	case ActionEndTurn:
		return v.validateEndTurn(action)
	case ActionFlipWildcard:
		return v.validateFlipWildcard(action)
	default:
		return Invalid(fmt.Sprintf("unknown action type: %s", action.Type))
	}
}

func (v *Validator) validateDrawCards(action *Action) ValidationResult {
	if v.state.TurnPhase != TurnPhaseDraw {
		return Invalid("can only draw cards at start of turn")
	}
	if action.PlayerID != v.state.GetCurrentPlayer().ID {
		return Invalid("not your turn")
	}
	return Valid()
}

func (v *Validator) validatePlayProperty(action *Action) ValidationResult {
	if v.state.TurnPhase != TurnPhaseActions {
		return Invalid("can only play cards during action phase")
	}
	if action.PlayerID != v.state.GetCurrentPlayer().ID {
		return Invalid("not your turn")
	}

	player := v.state.GetPlayer(action.PlayerID)
	card := player.GetCardFromHand(action.GetCardID())
	if card == nil {
		return Invalid("card not in hand")
	}

	cardType := card.GetType()
	if cardType != CardTypeProperty && cardType != CardTypeWildcard {
		return Invalid("card is not a property")
	}

	if v.state.ActionsPlayedThisTurn >= v.state.MaxActionsPerTurn {
		return Invalid("maximum actions per turn exceeded")
	}

	return Valid()
}

func (v *Validator) validatePlayMoney(action *Action) ValidationResult {
	if v.state.TurnPhase != TurnPhaseActions {
		return Invalid("can only play cards during action phase")
	}
	if action.PlayerID != v.state.GetCurrentPlayer().ID {
		return Invalid("not your turn")
	}

	player := v.state.GetPlayer(action.PlayerID)
	card := player.GetCardFromHand(action.GetCardID())
	if card == nil {
		return Invalid("card not in hand")
	}

	if card.GetValue() <= 0 {
		return Invalid("card has no money value")
	}

	if v.state.ActionsPlayedThisTurn >= v.state.MaxActionsPerTurn {
		return Invalid("maximum actions per turn exceeded")
	}

	return Valid()
}

func (v *Validator) validatePlayRent(action *Action) ValidationResult {
	if v.state.TurnPhase != TurnPhaseActions {
		return Invalid("can only play rent during action phase")
	}
	if action.PlayerID != v.state.GetCurrentPlayer().ID {
		return Invalid("not your turn")
	}

	player := v.state.GetPlayer(action.PlayerID)
	card := player.GetCardFromHand(action.GetCardID())
	if card == nil {
		return Invalid("card not in hand")
	}

	rentCard, ok := card.(*RentCard)
	if !ok {
		return Invalid("not a rent card")
	}

	color := action.GetColor()
	propertySet := player.GetPropertySet(color)

	// Must own at least one property of that color (unless wild rent)
	if len(propertySet.Cards) == 0 && !rentCard.IsWildRent() {
		return Invalid(fmt.Sprintf("you don't own any %s properties", color))
	}

	if v.state.ActionsPlayedThisTurn >= v.state.MaxActionsPerTurn {
		return Invalid("maximum actions per turn exceeded")
	}

	return Valid()
}

func (v *Validator) validatePlayAction(action *Action) ValidationResult {
	if v.state.TurnPhase != TurnPhaseActions {
		return Invalid("can only play actions during action phase")
	}
	if action.PlayerID != v.state.GetCurrentPlayer().ID {
		return Invalid("not your turn")
	}

	if v.state.ActionsPlayedThisTurn >= v.state.MaxActionsPerTurn {
		return Invalid("maximum actions per turn exceeded")
	}

	return Valid()
}

func (v *Validator) validateEndTurn(action *Action) ValidationResult {
	if action.PlayerID != v.state.GetCurrentPlayer().ID {
		return Invalid("not your turn")
	}
	return Valid()
}

func (v *Validator) validateFlipWildcard(action *Action) ValidationResult {
	if action.PlayerID != v.state.GetCurrentPlayer().ID {
		return Invalid("not your turn")
	}

	player := v.state.GetPlayer(action.PlayerID)
	cardID := action.GetCardID()

	// Find the wildcard in player's properties
	for _, ps := range player.Properties {
		for _, card := range ps.Cards {
			if card.GetID() == cardID {
				if _, ok := card.(*WildcardCard); ok {
					// Can't flip if it would break a complete set
					if ps.IsComplete() {
						return Invalid("cannot flip wildcard in complete set")
					}
					return Valid()
				}
				return Invalid("card is not a wildcard")
			}
		}
	}

	return Invalid("wildcard not found in properties")
}

func (v *Validator) validateResponse(action *Action) ValidationResult {
	if v.state.TurnPhase != TurnPhaseResponse {
		return Invalid("no pending action to respond to")
	}

	if !v.state.HasPendingAction() {
		return Invalid("no pending action")
	}

	// Check if this player is a target who hasn't responded
	isTarget := false
	for _, targetID := range v.state.GetPendingTargets() {
		if targetID == action.PlayerID {
			isTarget = true
			break
		}
	}
	if !isTarget {
		return Invalid("you are not a target of this action")
	}

	// For Fare Evasion, check they have the card
	if action.Type == ActionPlayFareEvasion {
		player := v.state.GetPlayer(action.PlayerID)
		card := player.GetCardFromHand(action.GetCardID())
		if card == nil {
			return Invalid("Fare Evasion card not in hand")
		}
		if card.GetName() != "Fare Evasion" {
			return Invalid("not a Fare Evasion card")
		}
	}

	return Valid()
}
